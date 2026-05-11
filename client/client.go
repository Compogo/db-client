package client

import (
	"context"
	"database/sql"
	"io"
)

// Client defines the interface for database operations.
// It mirrors the standard database/sql package while allowing
// for different driver implementations and decorators.
type Client interface {
	io.Closer

	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
	Exec(string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	SQL() *sql.DB
	DriverName() string
}
