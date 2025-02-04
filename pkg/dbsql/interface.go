package dbsql

import (
	"context"
	"database/sql"
)

type SqlBackend interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type DBConnectOptions struct {
	VerboseLogging bool
	DB             SqlBackend
	Filename       string // i.e "db.sqlite3"
}
