package db_client

import (
	"github.com/Compogo/compogo/container"
	"github.com/Compogo/db-client/client"
	"github.com/Compogo/types/linker"
)

var (
	// drivers stores constructor functions for each registered driver.
	// The linker associates each DriverName with its corresponding Getter.
	drivers = linker.NewLinker[string, Getter](linker.KeyStringNormalizer[Getter]())
)

// Registration registers a new database driver and its constructor function.
// This function should be called during driver package initialization.
// Once registered, the driver becomes available for use via the --db.driver flag.
func Registration(driverName string, getter Getter) {
	drivers.Add(driverName, getter)
}

// Getter is a function type that creates a new database client instance.
// It receives the DI container which may contain dependencies like config or logger,
// and returns a configured Client or an error.
type Getter func(container container.Container) (client.Client, error)
