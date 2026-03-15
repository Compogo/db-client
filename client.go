package db_client

import (
	"context"
	"database/sql"
	"fmt"
	"io"

	"github.com/Compogo/compogo/container"
	"github.com/Compogo/compogo/logger"
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
	Driver() Driver
}

// NewClient creates a database client instance based on the configured driver.
// It looks up the Getter for the selected driver, invokes it with the container,
// and returns the created client. Returns an error if the driver or its getter
// is not found, or if client creation fails.
func NewClient(config *Config, container container.Container, logger logger.Logger) (Client, error) {
	getter, err := getters.Get(config.Driver)
	if err != nil {
		return nil, fmt.Errorf("[db-client] driver '%s' getter undefined: %w", config.Driver, err)
	}

	client, err := getter(container)
	if err != nil {
		return nil, fmt.Errorf("[db-client] driver '%s' create failed: %w", config.Driver, err)
	}

	logger.Infof("[db-client] usage driver - '%s'", config.Driver)

	return client, nil
}
