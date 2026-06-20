package db_client

import (
	"fmt"
	"strings"

	"github.com/Compogo/compogo"
	"github.com/Compogo/compogo/flag"
)

// Component — компонент клиента БД для Compogo.
// Регистрирует конфигурацию и клиент в DI-контейнере.
//
// Пример подключения:
//
//	app.AddComponents(&db_client.Component)
//
//	var client db_client.Client
//	container.Invoke(func(c db_client.Client) { client = c })
//	rows, err := client.Query("SELECT * FROM users")
var Component = compogo.Component{
	Init: compogo.StepFunc(func(container compogo.Container) error {
		return container.Provides(
			NewConfig,
			NewClient,
		)
	}),
	BindFlags: compogo.BindFlags(func(flagSet flag.FlagSet, container compogo.Container) error {
		return container.Invoke(func(config *Config) {
			allDrivers := drivers.Keys()
			if len(allDrivers) == 1 {
				DriverDefault = allDrivers[0]
			}

			flagSet.StringVar(&config.Driver, DriverFieldName, DriverDefault, fmt.Sprintf("db client driver. Available drivers: [%s]", strings.Join(allDrivers, ",")))
		})
	}),
	Configuration: compogo.StepFunc(func(container compogo.Container) error {
		return container.Invoke(Configuration)
	}),
	Stop: compogo.StepFunc(func(container compogo.Container) error {
		return container.Invoke(func(c Client) error {
			return c.Close()
		})
	}),
}
