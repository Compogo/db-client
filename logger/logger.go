package logger

import (
	"context"
	"database/sql"

	"github.com/Compogo/compogo/logger"
	"github.com/Compogo/db-client/client"
)

// Logger decorates any dbClient.Client with query logging.
// All database operations are logged at DEBUG level, showing the query
// and its arguments. This is useful for debugging and monitoring.
type Logger struct {
	client.Client

	logger logger.Logger
}

func NewLogger(client client.Client, logger logger.Logger) *Logger {
	return &Logger{Client: client, logger: logger.GetLogger("db").GetLogger(client.DriverName())}
}

func (l *Logger) Close() error {
	l.logger.Info("close")

	return l.Client.Close()
}

func (l *Logger) Query(s string, i ...interface{}) (*sql.Rows, error) {
	l.logger.Debugf("query: %s; args: %+v", s, i)

	return l.Client.Query(s, i...)
}

func (l *Logger) QueryRow(s string, i ...interface{}) *sql.Row {
	l.logger.Debugf("query: %s; args: %+v", s, i)

	return l.Client.QueryRow(s, i...)
}

func (l *Logger) Exec(s string, i ...interface{}) (sql.Result, error) {
	l.logger.Debugf("query: %s; args: %+v", s, i)

	return l.Client.Exec(s, i...)
}

func (l *Logger) QueryContext(ctx context.Context, s string, i ...interface{}) (*sql.Rows, error) {
	l.logger.Debugf("query: %s; args: %+v", s, i)

	return l.Client.QueryContext(ctx, s, i...)
}

func (l *Logger) QueryRowContext(ctx context.Context, s string, i ...interface{}) *sql.Row {
	l.logger.Debugf("query: %s; args: %+v", s, i)

	return l.Client.QueryRowContext(ctx, s, i...)
}

func (l *Logger) ExecContext(ctx context.Context, s string, i ...interface{}) (sql.Result, error) {
	l.logger.Debugf("query: %s; args: %+v", s, i)

	return l.Client.ExecContext(ctx, s, i...)
}

func (l *Logger) SQL() *sql.DB {
	return l.Client.SQL()
}

func (l *Logger) DriverName() string {
	return l.Client.DriverName()
}
