package repository

import (
	"context"
	"database/sql"
)

type DBIImplRepository struct {
	db *sql.DB
}

func NewDBIImplRepository(db *sql.DB) *DBIImplRepository {
	return &DBIImplRepository{
		db: db,
	}
}

func (r DBIImplRepository) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, opts)
}

func (r DBIImplRepository) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return r.db.ExecContext(ctx, query, args...)
}

func (r DBIImplRepository) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return r.db.Query(query, args...)
}
