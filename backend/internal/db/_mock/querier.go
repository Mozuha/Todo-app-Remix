// Code generated by MockGen. DO NOT EDIT.
// Source: internal/db/wrapped_querier.go
//
// Generated by this command:
//
//	mockgen -source internal/db/wrapped_querier.go -destination internal/db/_mock/querier.go -package mock_db
//

// Package mock_db is a generated GoMock package.
package mock_db

import (
	context "context"
	reflect "reflect"
	db "todo-app/internal/db"

	pgx "github.com/jackc/pgx/v5"
	pgtype "github.com/jackc/pgx/v5/pgtype"
	gomock "go.uber.org/mock/gomock"
)

// MockWrappedQuerier is a mock of WrappedQuerier interface.
type MockWrappedQuerier struct {
	ctrl     *gomock.Controller
	recorder *MockWrappedQuerierMockRecorder
	isgomock struct{}
}

// MockWrappedQuerierMockRecorder is the mock recorder for MockWrappedQuerier.
type MockWrappedQuerierMockRecorder struct {
	mock *MockWrappedQuerier
}

// NewMockWrappedQuerier creates a new mock instance.
func NewMockWrappedQuerier(ctrl *gomock.Controller) *MockWrappedQuerier {
	mock := &MockWrappedQuerier{ctrl: ctrl}
	mock.recorder = &MockWrappedQuerierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWrappedQuerier) EXPECT() *MockWrappedQuerierMockRecorder {
	return m.recorder
}

// CreateTodo mocks base method.
func (m *MockWrappedQuerier) CreateTodo(ctx context.Context, arg db.CreateTodoParams) (db.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTodo", ctx, arg)
	ret0, _ := ret[0].(db.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTodo indicates an expected call of CreateTodo.
func (mr *MockWrappedQuerierMockRecorder) CreateTodo(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTodo", reflect.TypeOf((*MockWrappedQuerier)(nil).CreateTodo), ctx, arg)
}

// CreateUser mocks base method.
func (m *MockWrappedQuerier) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, arg)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockWrappedQuerierMockRecorder) CreateUser(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockWrappedQuerier)(nil).CreateUser), ctx, arg)
}

// DeleteTodo mocks base method.
func (m *MockWrappedQuerier) DeleteTodo(ctx context.Context, arg db.DeleteTodoParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTodo", ctx, arg)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTodo indicates an expected call of DeleteTodo.
func (mr *MockWrappedQuerierMockRecorder) DeleteTodo(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTodo", reflect.TypeOf((*MockWrappedQuerier)(nil).DeleteTodo), ctx, arg)
}

// DeleteUser mocks base method.
func (m *MockWrappedQuerier) DeleteUser(ctx context.Context, userID pgtype.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockWrappedQuerierMockRecorder) DeleteUser(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockWrappedQuerier)(nil).DeleteUser), ctx, userID)
}

// GetUserByEmail mocks base method.
func (m *MockWrappedQuerier) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", ctx, email)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockWrappedQuerierMockRecorder) GetUserByEmail(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockWrappedQuerier)(nil).GetUserByEmail), ctx, email)
}

// GetUserByID mocks base method.
func (m *MockWrappedQuerier) GetUserByID(ctx context.Context, id int32) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", ctx, id)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockWrappedQuerierMockRecorder) GetUserByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockWrappedQuerier)(nil).GetUserByID), ctx, id)
}

// GetUserByUserID mocks base method.
func (m *MockWrappedQuerier) GetUserByUserID(ctx context.Context, userID pgtype.UUID) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUserID", ctx, userID)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUserID indicates an expected call of GetUserByUserID.
func (mr *MockWrappedQuerierMockRecorder) GetUserByUserID(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUserID", reflect.TypeOf((*MockWrappedQuerier)(nil).GetUserByUserID), ctx, userID)
}

// ListTodos mocks base method.
func (m *MockWrappedQuerier) ListTodos(ctx context.Context, userID int32) ([]db.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTodos", ctx, userID)
	ret0, _ := ret[0].([]db.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTodos indicates an expected call of ListTodos.
func (mr *MockWrappedQuerierMockRecorder) ListTodos(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTodos", reflect.TypeOf((*MockWrappedQuerier)(nil).ListTodos), ctx, userID)
}

// SearchTodos mocks base method.
func (m *MockWrappedQuerier) SearchTodos(ctx context.Context, arg db.SearchTodosParams) ([]db.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchTodos", ctx, arg)
	ret0, _ := ret[0].([]db.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchTodos indicates an expected call of SearchTodos.
func (mr *MockWrappedQuerierMockRecorder) SearchTodos(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchTodos", reflect.TypeOf((*MockWrappedQuerier)(nil).SearchTodos), ctx, arg)
}

// UpdateTodo mocks base method.
func (m *MockWrappedQuerier) UpdateTodo(ctx context.Context, arg db.UpdateTodoParams) (db.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTodo", ctx, arg)
	ret0, _ := ret[0].(db.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTodo indicates an expected call of UpdateTodo.
func (mr *MockWrappedQuerierMockRecorder) UpdateTodo(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTodo", reflect.TypeOf((*MockWrappedQuerier)(nil).UpdateTodo), ctx, arg)
}

// UpdateTodoPosition mocks base method.
func (m *MockWrappedQuerier) UpdateTodoPosition(ctx context.Context, arg db.UpdateTodoPositionParams) (db.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTodoPosition", ctx, arg)
	ret0, _ := ret[0].(db.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTodoPosition indicates an expected call of UpdateTodoPosition.
func (mr *MockWrappedQuerierMockRecorder) UpdateTodoPosition(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTodoPosition", reflect.TypeOf((*MockWrappedQuerier)(nil).UpdateTodoPosition), ctx, arg)
}

// UpdateUsername mocks base method.
func (m *MockWrappedQuerier) UpdateUsername(ctx context.Context, arg db.UpdateUsernameParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUsername", ctx, arg)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUsername indicates an expected call of UpdateUsername.
func (mr *MockWrappedQuerierMockRecorder) UpdateUsername(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUsername", reflect.TypeOf((*MockWrappedQuerier)(nil).UpdateUsername), ctx, arg)
}

// WithTx mocks base method.
func (m *MockWrappedQuerier) WithTx(tx pgx.Tx) db.WrappedQuerier {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithTx", tx)
	ret0, _ := ret[0].(db.WrappedQuerier)
	return ret0
}

// WithTx indicates an expected call of WithTx.
func (mr *MockWrappedQuerierMockRecorder) WithTx(tx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTx", reflect.TypeOf((*MockWrappedQuerier)(nil).WithTx), tx)
}