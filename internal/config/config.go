package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

// Config stores the configuration for the application.
type Config struct {
	Port int64
}

// ParseFromEnvs parses the configuration from environment variables.
func ParseFromEnvs() (*Config, error) {
	var config Config
	var errs error

	if portStr := os.Getenv("ESPROFILER_PORT"); portStr != "" {
		var err error
		config.Port, err = strconv.ParseInt(portStr, 10, 64)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("failed to parse ESPROFILER_PORT: %w", err))
		}
	}

	if errs != nil {
		return nil, errs
	}
	return &config, nil
}
