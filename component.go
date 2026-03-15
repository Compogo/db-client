package db_client

import (
	"fmt"
	"strings"

	"github.com/Compogo/compogo/component"
	"github.com/Compogo/compogo/container"
	"github.com/Compogo/compogo/flag"
	"github.com/Compogo/db-client/client"
)

// Component is a ready-to-use Compogo component that provides a database client.
// It automatically:
//   - Registers Config and client factory in the DI container
//   - Adds command-line flags for driver selection
//   - Discovers available drivers from the registry
//   - Creates the appropriate client based on the selected driver
//
// Usage:
//
//	compogo.WithComponents(
//	    db_client.Component,
//	    // ... driver components (postgres, mysql, etc.)
//	)
//
// The actual client instance can be injected into any component that needs it:
//
//	type UserService struct {
//	    db db_client.Client
//	}
var Component = &component.Component{
	Init: component.StepFunc(func(container container.Container) error {
		return container.Provides(
			NewConfig,
			NewClient,
			func(client client.Client) Client { return client },
		)
	}),
	BindFlags: component.BindFlags(func(flagSet flag.FlagSet, container container.Container) error {
		return container.Invoke(func(config *Config) {
			allDrivers := drivers.Keys()
			if len(allDrivers) == 1 {
				DriverDefault = allDrivers[0]
			}

			flagSet.StringVar(&config.driver, DriverFieldName, DriverDefault, fmt.Sprintf("db client driver. Available drivers: [%s]", strings.Join(allDrivers, ",")))
		})
	}),
	Configuration: component.StepFunc(func(container container.Container) error {
		return container.Invoke(Configuration)
	}),
}
