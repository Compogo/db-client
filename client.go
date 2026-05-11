package db_client

import (
	"fmt"

	"github.com/Compogo/compogo/container"
	"github.com/Compogo/compogo/logger"
	"github.com/Compogo/db-client/client"
)

type Client client.Client

// NewClient creates a database client instance based on the configured driver.
// It looks up the Getter for the selected driver, invokes it with the container,
// and returns the created client. Returns an error if the driver or its getter
// is not found, or if client creation fails.
func NewClient(config *Config, container container.Container, logger logger.Logger) (client.Client, error) {
	getter, err := drivers.Get(config.Driver)
	if err != nil {
		return nil, fmt.Errorf("[db-client] driver '%s' getter undefined: %w", config.Driver, err)
	}

	c, err := getter(container)
	if err != nil {
		return nil, fmt.Errorf("[db-client] driver '%s' create failed: %w", config.Driver, err)
	}

	logger.Infof("[db-client] usage driver - '%s'", config.Driver)

	return c, nil
}
