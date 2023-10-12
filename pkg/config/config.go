package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServePort int    `envconfig:"SERVE_PORT" required:"true"`
	LogLevel  string `envconfig:"LOG_LEVEL" required:"false"`
}

// FromENV loads the environment variables to Config.
func FromENV() (*Config, error) {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}

	return &config, err
}

func (c *Config) AppAddress() string {
	return fmt.Sprintf("0.0.0.0:%d", c.ServePort)
}

func (c *Config) GetLogLevel() string {
	return c.LogLevel
}
