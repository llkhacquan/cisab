package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	// Server configuration
	Server ServerConfig `yaml:"server"`

	// Environment (dev, staging, production)
	Environment string `yaml:"environment"`
}

// ServerConfig holds the server-related configuration
type ServerConfig struct {
	// Port to run the server on
	Port int `yaml:"port"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: 8080,
		},
		Environment: "dev",
	}
}

// LoadConfig loads the configuration from the specified file path
func LoadConfig(path string) (*Config, error) {
	// Set default config
	config := DefaultConfig()

	// Read the config file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "error reading config file")
	}

	// Unmarshal the config file
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling config")
	}

	// Override with environment variables if they exist
	if port := os.Getenv("PORT"); port != "" {
		var p int
		if _, err := fmt.Sscanf(port, "%d", &p); err == nil {
			config.Server.Port = p
		}
	}

	if env := os.Getenv("ENVIRONMENT"); env != "" {
		config.Environment = env
	}

	return config, nil
}

// LoadConfigFromEnv loads the configuration based on the current environment
func LoadConfigFromEnv() (*Config, error) {
	// Default to local environment if not specified
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "local"
	}

	// Determine config file path
	configPath := filepath.Join("config", fmt.Sprintf("%s.yaml", env))

	// Try to load the config file
	config, err := LoadConfig(configPath)
	if err != nil {
		// If we can't find the config file, try loading from the default path
		defaultPath := filepath.Join("config", "local.yaml")
		if env != "local" {
			config, err = LoadConfig(defaultPath)
			if err != nil {
				return DefaultConfig(), err
			}
			return config, nil
		}
		return DefaultConfig(), err
	}

	return config, nil
}
