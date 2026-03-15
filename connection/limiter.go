package connection

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Compogo/db-client/client"
	driver2 "github.com/Compogo/db-client/driver"
	"github.com/Compogo/types/counter"
)

// Limiter implements a circuit breaker pattern for database connections.
// It tracks consecutive errors and blocks requests for a duration when
// the error limit is reached. After the duration, it allows requests again.
type Limiter struct {
	client.Client

	counter counter.Counter
	limit   int64

	timer    *time.Timer
	duration time.Duration
	m        sync.Mutex
}

func NewLimiter(client client.Client, limit int64, duration time.Duration) *Limiter {
	return &Limiter{Client: client, limit: limit, duration: duration, counter: counter.NewCounter()}
}

func (l *Limiter) Close() error {
	defer l.counter.Reset()
	defer func() {
		l.m.Lock()
		defer l.m.Unlock()

		if l.timer != nil {
			l.timer.Stop()
			l.timer = nil
		}
	}()

	return l.Client.Close()
}

func (l *Limiter) Query(s string, i ...interface{}) (*sql.Rows, error) {
	if l.counter.Get() >= l.limit {
		if err := l.timeoutProcess(); err != nil {
			return nil, err
		}
	}

	rows, err := l.Client.Query(s, i...)
	if errors.Is(err, sql.ErrConnDone) || errors.Is(err, driver.ErrBadConn) {
		l.counter.Inc()
		return nil, err
	}

	l.counter.Reset()

	return rows, err
}

func (l *Limiter) QueryRow(s string, i ...interface{}) *sql.Row {
	rows := l.Client.QueryRow(s, i...)

	if rows != nil && (errors.Is(rows.Err(), sql.ErrConnDone) || errors.Is(rows.Err(), driver.ErrBadConn)) {
		l.counter.Inc()
	} else {
		l.counter.Reset()
	}

	return rows
}

func (l *Limiter) Exec(s string, i ...interface{}) (sql.Result, error) {
	if l.counter.Get() >= l.limit {
		if err := l.timeoutProcess(); err != nil {
			return nil, err
		}
	}

	result, err := l.Client.Exec(s, i...)
	if errors.Is(err, sql.ErrConnDone) || errors.Is(err, driver.ErrBadConn) {
		l.counter.Inc()
		return nil, err
	}

	l.counter.Reset()

	return result, err
}

func (l *Limiter) QueryContext(ctx context.Context, s string, i ...interface{}) (*sql.Rows, error) {
	if l.counter.Get() >= l.limit {
		if err := l.timeoutProcess(); err != nil {
			return nil, err
		}
	}

	rows, err := l.Client.QueryContext(ctx, s, i...)
	if errors.Is(err, sql.ErrConnDone) || errors.Is(err, driver.ErrBadConn) {
		l.counter.Inc()
		return nil, err
	}

	l.counter.Reset()

	return rows, err
}

func (l *Limiter) QueryRowContext(ctx context.Context, s string, i ...interface{}) *sql.Row {
	rows := l.Client.QueryRowContext(ctx, s, i...)

	if rows != nil && (errors.Is(rows.Err(), sql.ErrConnDone) || errors.Is(rows.Err(), driver.ErrBadConn)) {
		l.counter.Inc()
	} else {
		l.counter.Reset()
	}

	return rows
}

func (l *Limiter) ExecContext(ctx context.Context, s string, i ...interface{}) (sql.Result, error) {
	if l.counter.Get() >= l.limit {
		if err := l.timeoutProcess(); err != nil {
			return nil, err
		}
	}

	result, err := l.Client.ExecContext(ctx, s, i...)
	if errors.Is(err, sql.ErrConnDone) || errors.Is(err, driver.ErrBadConn) {
		l.counter.Inc()
		return nil, err
	}

	l.counter.Reset()

	return result, err
}

func (l *Limiter) Driver() driver2.Driver {
	return l.Client.Driver()
}

func (l *Limiter) SQL() *sql.DB {
	return l.Client.SQL()
}

func (l *Limiter) timeoutProcess() (err error) {
	l.m.Lock()
	defer l.m.Unlock()

	if l.timer == nil {
		l.timer = time.NewTimer(l.duration)
	}

	select {
	case <-l.timer.C:
		l.timer = nil
		return nil
	default:
		return fmt.Errorf("[db-client.limiter] driver '%s' error: %w", l.Driver(), driver.ErrBadConn)
	}
}
