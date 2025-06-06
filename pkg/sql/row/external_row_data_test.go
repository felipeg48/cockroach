// Copyright 2024 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package row_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/sql"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descs"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/testutils/serverutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/skip"
	"github.com/cockroachdb/cockroach/pkg/testutils/sqlutils"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
	"github.com/stretchr/testify/require"
)

// TestExternalRowData is a sanity test that external row data (as configured by
// the External field of the table descriptor) is accessed correctly. It does so
// by creating two tables with one pointing to the other at a specific point in
// time.
func TestExternalRowData(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	ctx := context.Background()
	srv, sqlDB, _ := serverutils.StartServer(t, base.TestServerArgs{})
	defer srv.Stopper().Stop(ctx)
	s := srv.ApplicationLayer()

	// Ensure that we always get the same connection in the SQL runner.
	sqlDB.SetMaxOpenConns(1)

	r := sqlutils.MakeSQLRunner(sqlDB)
	r.Exec(t, `CREATE TABLE t (k INT PRIMARY KEY, v1 INT, v2 INT, INDEX (v1))`)
	r.Exec(t, `CREATE TABLE t_copy (k INT PRIMARY KEY, v1 INT, v2 INT, INDEX (v1))`)

	// Insert some data into the original table, then record AOST, and insert
	// more data that shouldn't be visible via the external copy.
	r.Exec(t, `INSERT INTO t SELECT i, i, -i FROM generate_series(1, 3) AS g(i)`)
	asOf := hlc.Timestamp{WallTime: timeutil.Now().UnixNano()}
	r.Exec(t, `INSERT INTO t SELECT i, i, -i FROM generate_series(4, 6) AS g(i)`)

	// Modify the table descriptor for 't_copy' to have external row data from
	// 't'.
	var tableID int
	row := r.QueryRow(t, `SELECT 't'::REGCLASS::OID`)
	row.Scan(&tableID)
	execCfg := s.ExecutorConfig().(sql.ExecutorConfig)
	require.NoError(t, execCfg.InternalDB.DescsTxn(ctx, func(ctx context.Context, txn descs.Txn) error {
		descriptors := txn.Descriptors()
		tn := tree.MakeTableNameWithSchema("defaultdb", "public", "t_copy")
		_, mut, err := descs.PrefixAndMutableTable(ctx, descriptors.MutableByName(txn.KV()), &tn)
		if err != nil {
			return err
		}
		require.NotNil(t, mut)
		mut.External = &descpb.ExternalRowData{
			AsOf:     asOf,
			TenantID: execCfg.Codec.TenantID,
			TableID:  descpb.ID(tableID),
		}
		return descriptors.WriteDesc(ctx, false /* kvTrace */, mut, txn.KV())
	}))

	// Try both execution engines since they have different fetcher
	// implementations.
	for _, vectorize := range []string{"on", "off"} {
		r.Exec(t, `SET vectorize = `+vectorize)
		for _, tc := range []struct {
			query    string
			expected [][]string
		}{
			{ // ScanRequest
				query:    `SELECT * FROM t_copy`,
				expected: [][]string{{"1", "1", "-1"}, {"2", "2", "-2"}, {"3", "3", "-3"}},
			},
			{ // ReverseScanRequest
				query:    `SELECT * FROM t_copy ORDER BY k DESC`,
				expected: [][]string{{"3", "3", "-3"}, {"2", "2", "-2"}, {"1", "1", "-1"}},
			},
			{ // GetRequests
				query:    `SELECT * FROM t_copy WHERE k = 2 OR k = 5`,
				expected: [][]string{{"2", "2", "-2"}},
			},
			{ // lookup join which might be served via the Streamer
				query:    `SELECT t_copy.k FROM t INNER LOOKUP JOIN t_copy ON t.k = t_copy.k`,
				expected: [][]string{{"1"}, {"2"}, {"3"}},
			},
			{ // index join which might be served via the Streamer
				query:    `SELECT * FROM t_copy WHERE v1 = 2 OR v1 = 5`,
				expected: [][]string{{"2", "2", "-2"}},
			},
		} {

			require.Equal(t, tc.expected, r.QueryStrMeta(
				t, fmt.Sprintf("vectorize=%v", vectorize), tc.query,
			))
		}
	}
}

// TestExternalRowDataDistSQL tests that the DistSQL physical planner can
// correctly place flows reading from external row data.
func TestExternalRowDataDistSQL(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	skip.UnderDuress(t, "slow test")

	ctx := context.Background()

	// Start a 5-node cluster.
	testCluster := serverutils.StartCluster(t, 5, /* numNodes */
		base.TestClusterArgs{
			ReplicationMode: base.ReplicationManual,
			ServerArgs: base.TestServerArgs{
				UseDatabase:       "defaultdb",
				DefaultTestTenant: base.TestControlsTenantsExplicitly,
			},
		})
	defer testCluster.Stopper().Stop(ctx)
	ts := testCluster.Server(0)

	srcTenant, srcDB, err := ts.TenantController().StartSharedProcessTenant(ctx,
		base.TestSharedProcessTenantArgs{
			TenantID:    serverutils.TestTenantID(),
			TenantName:  "src",
			UseDatabase: "defaultdb",
		},
	)
	require.NoError(t, err)
	dstTenant, dstDB, err := ts.TenantController().StartSharedProcessTenant(ctx,
		base.TestSharedProcessTenantArgs{
			TenantID:    serverutils.TestTenantID2(),
			TenantName:  "dst",
			UseDatabase: "defaultdb",
		},
	)
	require.NoError(t, err)

	// Set up the source table.
	srcRunner := sqlutils.MakeSQLRunner(srcDB)
	otherRunner := sqlutils.MakeSQLRunner(dstDB)

	// Set up the source table. Place leaseholders on nodes 3, 4, 5.
	srcRunner.Exec(t, `CREATE TABLE t_src (k INT PRIMARY KEY, v1 INT, v2 INT)`)
	srcRunner.Exec(t, `INSERT INTO t_src VALUES (1), (3), (5)`)
	srcRunner.Exec(t, `ALTER TABLE t_src SPLIT AT VALUES (2), (4)`)
	srcRunner.ExecSucceedsSoon(
		t, `ALTER TABLE t_src RELOCATE VALUES (ARRAY[3], 1), (ARRAY[4], 3), (ARRAY[5], 5)`,
	)

	for _, tc := range []struct {
		name          string
		dstRunner     *sqlutils.SQLRunner
		dstInternalDB *sql.InternalDB
	}{
		{
			name:          "same-tenant",
			dstRunner:     srcRunner,
			dstInternalDB: srcTenant.InternalDB().(*sql.InternalDB),
		},
		{
			name:          "different-tenant",
			dstRunner:     otherRunner,
			dstInternalDB: dstTenant.InternalDB().(*sql.InternalDB),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.dstRunner.Exec(t, `CREATE TABLE t_dst (k INT PRIMARY KEY, v1 INT, v2 INT)`)
			asOf := hlc.Timestamp{WallTime: timeutil.Now().UnixNano()}

			// Modify the table descriptor for 't_dst' to have external row data from
			// table 't_src'.
			var tableID int
			row := srcRunner.QueryRow(t, `SELECT 't_src'::REGCLASS::OID`)
			row.Scan(&tableID)
			require.NoError(t, tc.dstInternalDB.DescsTxn(ctx, func(ctx context.Context, txn descs.Txn) error {
				descriptors := txn.Descriptors()
				tn := tree.MakeTableNameWithSchema("defaultdb", "public", "t_dst")
				_, mut, err := descs.PrefixAndMutableTable(ctx, descriptors.MutableByName(txn.KV()), &tn)
				if err != nil {
					return err
				}
				require.NotNil(t, mut)
				mut.External = &descpb.ExternalRowData{
					AsOf:     asOf,
					TenantID: serverutils.TestTenantID(),
					TableID:  descpb.ID(tableID),
				}
				return descriptors.WriteDesc(ctx, false /* kvtrace */, mut, txn.KV())
			}))

			// Now check that DistSQL plans against both tables correctly place
			// flows on nodes 1, 3, 4, 5.
			srcRunner.Exec(t, `SET distsql = always`)
			tc.dstRunner.Exec(t, `SET distsql = always`)
			exp := `"nodeNames":["1","3","4","5"]`
			var info string
			row = srcRunner.QueryRow(t, `EXPLAIN (DISTSQL, JSON) SELECT count(*) FROM t_src`)
			row.Scan(&info)
			if !strings.Contains(info, exp) {
				t.Fatalf("expected DistSQL plan to contain %s: was %s", exp, info)
			}
			row = tc.dstRunner.QueryRow(t, `EXPLAIN (DISTSQL, JSON) SELECT count(*) FROM t_dst`)
			row.Scan(&info)
			if !strings.Contains(info, exp) {
				t.Fatalf("expected DistSQL plan to contain %s: was %s", exp, info)
			}
		})
	}
}
