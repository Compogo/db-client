package db_client

import (
	"errors"
	"fmt"

	"github.com/Compogo/compogo/configurator"
	"github.com/Compogo/db-client/driver"
)

const (
	DriverFieldName = "db.driver"
)

var (
	DriverDefault = ""
)

type Config struct {
	driver string
	Driver driver.Driver
}

func NewConfig() *Config {
	return &Config{}
}

func Configuration(config *Config, configurator configurator.Configurator) (*Config, error) {
	if config.driver == "" || config.driver == DriverDefault {
		configurator.SetDefault(DriverFieldName, DriverDefault)
		config.driver = configurator.GetString(DriverFieldName)
	}

	if config.driver == "" {
		return nil, errors.New("[db.client] driver is not set")
	}

	d, err := drivers.Get(config.driver)
	if err != nil {
		return nil, fmt.Errorf("[db.client] driver '%s' get failed: %w", config.driver, err)
	}

	config.Driver = d

	return config, nil
}
