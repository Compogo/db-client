package db_client

import (
	"errors"
	"fmt"

	"github.com/Compogo/compogo"
)

// DriverFieldName — имя поля для драйвера БД.
const DriverFieldName = "db.driver"

// DriverDefault — драйвер по умолчанию (пустая строка, требует явного указания).
var DriverDefault = ""

// Config содержит конфигурацию клиента БД.
type Config struct {
	Driver string
}

// NewConfig создаёт новую конфигурацию.
func NewConfig() *Config {
	return &Config{}
}

// Configuration загружает конфигурацию из Configurator.
// Проверяет, что драйвер указан и зарегистрирован.
func Configuration(config *Config, configurator compogo.Configurator) (*Config, error) {
	if config.Driver == "" || config.Driver == DriverDefault {
		configurator.SetDefault(DriverFieldName, DriverDefault)
		config.Driver = configurator.GetString(DriverFieldName)
	}

	if config.Driver == "" {
		return nil, errors.New("[Database] driver is not set")
	}

	if !drivers.Has(config.Driver) {
		return nil, fmt.Errorf("[Database] driver '%s' unknown", config.Driver)
	}

	return config, nil
}
