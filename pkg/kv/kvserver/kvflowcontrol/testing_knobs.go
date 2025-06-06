// Copyright 2023 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package kvflowcontrol

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/base"
)

// TestingKnobs provide fine-grained control over the various kvflowcontrol
// components for testing.
type TestingKnobs struct {
	UseOnlyForScratchRanges bool
	// OverrideTokenDeduction is used to override how many tokens are deducted
	// post-evaluation.
	OverrideTokenDeduction func(tokens Tokens) Tokens
	// OverridePullPushMode is used to override whether the pull mode, or push
	// mode is enabled.
	//
	// - when set to true, pull mode is enabled
	// - when set to false, push mode is enabled
	// - when left unset the otherwise set mode is used
	//
	// This is used to test the behavior of the flow control in push and pull
	// mode, while also having the ability to switch between the two
	// apply_to_(elastic|all) modes.
	OverridePullPushMode func() bool
	// OverrideBypassAdmitWaitForEval is used to override the behavior of
	// WaitForEval. When bypass is set to true, WaitForEval will return
	// immediately and return the waited value. Otherwise, when bypass is set
	// to false, or unset, WaitForEval will behave normally.
	OverrideBypassAdmitWaitForEval func(ctx context.Context) (bypass bool, waited bool)
	// OverrideAlwaysRefreshSendStreamStats is used to override the behavior of
	// the send stream stats refresh. When set to true, the send stream stats
	// will always be refreshed on a HandleRaftEventRaftMuLocked call. Otherwise,
	// when set to false, the default behavior will be used.
	OverrideAlwaysRefreshSendStreamStats bool
}

// ModuleTestingKnobs is part of the base.ModuleTestingKnobs interface.
func (t *TestingKnobs) ModuleTestingKnobs() {}

var _ base.ModuleTestingKnobs = (*TestingKnobs)(nil)
