package db_client

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/Compogo/runner"
)

// Limiter — обёртка над Client с ограничением на количество ошибок соединения.
// При превышении лимита ошибок (соединений) переходит в режим ожидания.
//
// Используется для защиты БД от "штурма" недоступным соединением.
type Limiter struct {
	Client

	counter atomic.Uint64
	limit   uint64

	isConnError atomic.Bool

	duration time.Duration
	runner   runner.Runner
}

// NewLimiter создаёт новый Limiter.
func NewLimiter(client Client, limit uint64, duration time.Duration, runner runner.Runner) *Limiter {
	return &Limiter{Client: client, limit: limit, duration: duration, runner: runner}
}

func (l *Limiter) Close() error {
	defer l.counter.Store(0)

	return l.Client.Close()
}

func (l *Limiter) Query(s string, i ...interface{}) (*sql.Rows, error) {
	if err := l.validate(); err != nil {
		return nil, err
	}

	rows, err := l.Client.Query(s, i...)
	if errors.Is(err, sql.ErrConnDone) || errors.Is(err, driver.ErrBadConn) {
		l.counter.Add(1)
		return nil, err
	}

	l.counter.Store(0)

	return rows, err
}

func (l *Limiter) QueryRow(s string, i ...interface{}) *sql.Row {
	rows := l.Client.QueryRow(s, i...)

	if rows != nil && (errors.Is(rows.Err(), sql.ErrConnDone) || errors.Is(rows.Err(), driver.ErrBadConn)) {
		l.counter.Add(1)
	} else {
		l.counter.Store(0)
	}

	return rows
}

func (l *Limiter) Exec(s string, i ...interface{}) (sql.Result, error) {
	if err := l.validate(); err != nil {
		return nil, err
	}

	result, err := l.Client.Exec(s, i...)
	if errors.Is(err, sql.ErrConnDone) || errors.Is(err, driver.ErrBadConn) {
		l.counter.Add(1)
		return nil, err
	}

	l.counter.Store(0)

	return result, err
}

func (l *Limiter) QueryContext(ctx context.Context, s string, i ...interface{}) (*sql.Rows, error) {
	if err := l.validate(); err != nil {
		return nil, err
	}

	rows, err := l.Client.QueryContext(ctx, s, i...)
	if errors.Is(err, sql.ErrConnDone) || errors.Is(err, driver.ErrBadConn) {
		l.counter.Add(1)
		return nil, err
	}

	l.counter.Store(0)

	return rows, err
}

func (l *Limiter) QueryRowContext(ctx context.Context, s string, i ...interface{}) *sql.Row {
	rows := l.Client.QueryRowContext(ctx, s, i...)

	if rows != nil && (errors.Is(rows.Err(), sql.ErrConnDone) || errors.Is(rows.Err(), driver.ErrBadConn)) {
		l.counter.Add(1)
	} else {
		l.counter.Store(0)
	}

	return rows
}

func (l *Limiter) ExecContext(ctx context.Context, s string, i ...interface{}) (sql.Result, error) {
	if err := l.validate(); err != nil {
		return nil, err
	}

	result, err := l.Client.ExecContext(ctx, s, i...)
	if errors.Is(err, sql.ErrConnDone) || errors.Is(err, driver.ErrBadConn) {
		l.counter.Add(1)
		return nil, err
	}

	l.counter.Store(0)

	return result, err
}

func (l *Limiter) validate() error {
	if l.isConnError.Load() {
		return fmt.Errorf("[Database][limiter] driver '%s' error: %w", l.DriverName(), driver.ErrBadConn)
	}

	if l.counter.Load() >= l.limit && !l.runner.HasProcess(l) {
		return l.runner.RunProcess(l)
	}

	return nil
}

func (l *Limiter) Process(ctx context.Context) error {
	l.isConnError.Store(true)
	defer l.isConnError.Store(false)
	defer l.counter.Store(0)

	ctx, cancel := context.WithTimeout(ctx, l.duration)
	defer cancel()

	<-ctx.Done()

	return nil
}

func (l *Limiter) Name() string {
	return "Database.limiter"
}
