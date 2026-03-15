package db_client

import (
	"github.com/Compogo/compogo/container"
	"github.com/Compogo/db-client/client"
	"github.com/Compogo/db-client/driver"
	"github.com/Compogo/types/linker"
	"github.com/Compogo/types/mapper"
)

var (
	// drivers stores all registered database drivers by their string identifiers.
	drivers = mapper.NewMapper[driver.Driver]()

	// getters stores constructor functions for each registered driver.
	// The linker associates each Driver with its corresponding Getter.
	getters = linker.NewLinker[driver.Driver, Getter]()
)

// Registration registers a new database driver and its constructor function.
// This function should be called during driver package initialization.
// Once registered, the driver becomes available for use via the --db.driver flag.
func Registration(d driver.Driver, getter Getter) {
	drivers.Add(d)
	getters.Add(d, getter)
}

// Getter is a function type that creates a new database client instance.
// It receives the DI container which may contain dependencies like config or logger,
// and returns a configured Client or an error.
type Getter func(container container.Container) (client.Client, error)
