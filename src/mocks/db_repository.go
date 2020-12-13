package repository

import (
	"context"
	"database/sql"
	mock "github.com/stretchr/testify/mock"
)

type MockDBRepository struct {
	mock.Mock
}

func (m *MockDBRepository) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	call := m.Called(ctx, opts)
	res := call.Get(0)
	if res == nil {
		return nil, call.Error(1)
	}
	return res.(*sql.Tx), call.Error(1)
}

// GetListTeller :
func (m *MockDBRepository) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	call := m.Called(ctx, query, args)
	res := call.Get(0)
	if res == nil {
		return nil, call.Error(1)
	}
	return res.(sql.Result), call.Error(1)
}

// GetListAATR :
func (m *MockDBRepository) Query(query string, args ...interface{}) (*sql.Rows, error) {
	call := m.Called(query, args)
	res := call.Get(0)
	if res == nil {
		return nil, call.Error(1)
	}
	return res.(*sql.Rows), call.Error(1)
}
