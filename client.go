package db_client

import (
	"context"
	"database/sql"
	"fmt"
	"io"

	"github.com/Compogo/compogo"
)

// Client — интерфейс для работы с БД.
// Объединяет стандартные методы database/sql с поддержкой контекста.
//
// Поддерживает:
//   - Выполнение запросов (Query, QueryRow, Exec)
//   - Контекстные версии (QueryContext, QueryRowContext, ExecContext)
//   - Доступ к *sql.DB для низкоуровневых операций
//   - Graceful shutdown через io.Closer
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

// NewClient создаёт новый клиент БД с указанным драйвером.
// Использует зарегистрированный Getter для создания клиента.
func NewClient(config *Config, container compogo.Container, logger compogo.Logger) (Client, error) {
	getter, err := drivers.Get(config.Driver)
	if err != nil {
		return nil, fmt.Errorf("[Database] driver '%s' getter undefined: %w", config.Driver, err)
	}

	c, err := getter(container)
	if err != nil {
		return nil, fmt.Errorf("[Database] driver '%s' create failed: %w", config.Driver, err)
	}

	logger.GetLogger("Database").Infof("usage driver - '%s'", config.Driver)

	return c, nil
}
