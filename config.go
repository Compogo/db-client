package db_client

import (
	"errors"
	"fmt"

	"github.com/Compogo/compogo/configurator"
)

const (
	DriverFieldName = "db.driver"
)

var (
	DriverDefault = ""
)

type Config struct {
	Driver string
}

func NewConfig() *Config {
	return &Config{}
}

func Configuration(config *Config, configurator configurator.Configurator) (*Config, error) {
	if config.Driver == "" || config.Driver == DriverDefault {
		configurator.SetDefault(DriverFieldName, DriverDefault)
		config.Driver = configurator.GetString(DriverFieldName)
	}

	if config.Driver == "" {
		return nil, errors.New("[db.client] driver is not set")
	}

	if !drivers.Has(config.Driver) {
		return nil, fmt.Errorf("[db.client] driver '%s' unknown", config.Driver)
	}

	return config, nil
}
