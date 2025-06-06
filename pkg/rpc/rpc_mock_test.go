// Copyright 2025 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.
//
// TODO(#147193): restore auto-generated RPC mock tests once gomock has support
// for generics.
package rpc

import (
	context "context"
	reflect "reflect"

	kvpb "github.com/cockroachdb/cockroach/pkg/kv/kvpb"
	roachpb "github.com/cockroachdb/cockroach/pkg/roachpb"
	rpcpb "github.com/cockroachdb/cockroach/pkg/rpc/rpcpb"
	gomock "github.com/golang/mock/gomock"
	"google.golang.org/grpc"
)

// MockBatchStreamClient is a mock of BatchStreamClient interface.
type MockBatchStreamClient struct {
	ctrl     *gomock.Controller
	recorder *MockBatchStreamClientMockRecorder
}

// MockBatchStreamClientMockRecorder is the mock recorder for MockBatchStreamClient.
type MockBatchStreamClientMockRecorder struct {
	mock *MockBatchStreamClient
}

// NewMockBatchStreamClient creates a new mock instance.
func NewMockBatchStreamClient(ctrl *gomock.Controller) *MockBatchStreamClient {
	mock := &MockBatchStreamClient{ctrl: ctrl}
	mock.recorder = &MockBatchStreamClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBatchStreamClient) EXPECT() *MockBatchStreamClientMockRecorder {
	return m.recorder
}

// Recv mocks base method.
func (m *MockBatchStreamClient) Recv() (*kvpb.BatchResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*kvpb.BatchResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockBatchStreamClientMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockBatchStreamClient)(nil).Recv))
}

// Send mocks base method.
func (m *MockBatchStreamClient) Send(arg0 *kvpb.BatchRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockBatchStreamClientMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockBatchStreamClient)(nil).Send), arg0)
}

// MockDialbacker is a mock of Dialbacker interface.
type MockDialbacker struct {
	ctrl     *gomock.Controller
	recorder *MockDialbackerMockRecorder
}

// MockDialbackerMockRecorder is the mock recorder for MockDialbacker.
type MockDialbackerMockRecorder struct {
	mock *MockDialbacker
}

// NewMockDialbacker creates a new mock instance.
func NewMockDialbacker(ctrl *gomock.Controller) *MockDialbacker {
	mock := &MockDialbacker{ctrl: ctrl}
	mock.recorder = &MockDialbackerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDialbacker) EXPECT() *MockDialbackerMockRecorder {
	return m.recorder
}

// GRPCDialNode mocks base method.
func (m *MockDialbacker) GRPCDialNode(
	arg0 string, arg1 roachpb.NodeID, arg2 roachpb.Locality, arg3 rpcpb.ConnectionClass,
) *Connection[*grpc.ClientConn] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GRPCDialNode", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*Connection[*grpc.ClientConn])
	return ret0
}

// GRPCDialNode indicates an expected call of GRPCDialNode.
func (mr *MockDialbackerMockRecorder) GRPCDialNode(
	arg0, arg1, arg2, arg3 interface{},
) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GRPCDialNode", reflect.TypeOf((*MockDialbacker)(nil).GRPCDialNode), arg0, arg1, arg2, arg3)
}

// GRPCUnvalidatedDial mocks base method.
func (m *MockDialbacker) GRPCUnvalidatedDial(
	arg0 string, arg1 roachpb.Locality,
) *Connection[*grpc.ClientConn] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GRPCUnvalidatedDial", arg0, arg1)
	ret0, _ := ret[0].(*Connection[*grpc.ClientConn])
	return ret0
}

// GRPCUnvalidatedDial indicates an expected call of GRPCUnvalidatedDial.
func (mr *MockDialbackerMockRecorder) GRPCUnvalidatedDial(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GRPCUnvalidatedDial", reflect.TypeOf((*MockDialbacker)(nil).GRPCUnvalidatedDial), arg0, arg1)
}

// grpcDialRaw mocks base method.
func (m *MockDialbacker) grpcDialRaw(
	arg0 context.Context, arg1 string, arg2 rpcpb.ConnectionClass, arg3 ...grpc.DialOption,
) (*grpc.ClientConn, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "grpcDialRaw", varargs...)
	ret0, _ := ret[0].(*grpc.ClientConn)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// grpcDialRaw indicates an expected call of grpcDialRaw.
func (mr *MockDialbackerMockRecorder) grpcDialRaw(
	arg0, arg1, arg2 interface{}, arg3 ...interface{},
) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "grpcDialRaw", reflect.TypeOf((*MockDialbacker)(nil).grpcDialRaw), varargs...)
}

// wrapCtx mocks base method.
func (m *MockDialbacker) wrapCtx(
	arg0 context.Context, arg1 string, arg2 roachpb.NodeID, arg3 rpcpb.ConnectionClass,
) context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "wrapCtx", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// wrapCtx indicates an expected call of wrapCtx.
func (mr *MockDialbackerMockRecorder) wrapCtx(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "wrapCtx", reflect.TypeOf((*MockDialbacker)(nil).wrapCtx), arg0, arg1, arg2, arg3)
}
