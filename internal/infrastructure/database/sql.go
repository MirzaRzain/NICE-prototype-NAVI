package db

import (
	"context"
	"database/sql"
)

type DB interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Query
}

type Query interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
