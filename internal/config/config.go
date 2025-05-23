// Package config provides functionality for loading and accessing application configuration.
// It supports loading from environment variables and configuration files.
package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration parameters.
type Config struct {
	ServerPort int
	HmacSecret []byte
}

// Load returns the application configuration.
func Load() (*Config, error) {
	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil || port == 0 {
		port = 3000
	}

	hmacSecret := []byte(os.Getenv("HMAC_SECRET"))

	return &Config{
		ServerPort: port,
		HmacSecret: hmacSecret,
	}, nil
}
