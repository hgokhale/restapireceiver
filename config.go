package restapireceiver

import (
	"fmt"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
	"strings"
)

type Config struct {
	scraperhelper.ControllerConfig `mapstructure:",squash"`
	Endpoint                       string `mapstructure:"endpoint"`
	AuthToken                      string `mapstructure:"auth_token"`
	Username                       string `mapstructure:"username"`
	Password                       string `mapstructure:"password"`
}

func (c *Config) Validate() error {
	var validationErrors []string = []string{}

	if c.Endpoint == "" {
		validationErrors = append(validationErrors, "'endpoint' is required")
	}

	if c.AuthToken == "" {
		if c.Username == "" || c.Password == "" {
			validationErrors = append(validationErrors, "either of 'auth_token' or 'username'+'password' are required")
		}
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("Config validation failed: %v", strings.Join(validationErrors, ", "))
	}
	return nil
}
