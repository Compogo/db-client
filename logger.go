package db_client

import (
	"context"
	"database/sql"

	"github.com/Compogo/compogo"
)

// Logger — обёртка над Client с логированием всех запросов.
// Логирует на уровне Debug все SQL-запросы и их аргументы.
type Logger struct {
	Client

	logger compogo.Logger
}

// NewLogger создаёт новый Logger для клиента БД.
func NewLogger(client Client, logger compogo.Logger) *Logger {
	return &Logger{Client: client, logger: logger.GetLogger("Database").GetLogger(client.DriverName())}
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
