package config

import (
	"fmt"
	"os"
)

const (
	//DefaultLogLevel is the default log level
	DefaultLogLevel = "info"
)

const (
	// OWMKeyEnvKey is the log level environment key
	OWMKeyEnvKey = "RHCC_OWM_API_KEY"

	// LogLevelEnvKey is the log level environment key
	LogLevelEnvKey = "RHCC_LOG_LEVEL"
)

// Config represents configuration for the service.
type Config struct {
	OWMKey     string
	LogLevel   string
	OutputJSON bool
}

// LoadConfig loads configuration from the environment.
func LoadConfig() (*Config, error) {
	owmKey := os.Getenv(OWMKeyEnvKey)
	if owmKey == "" {
		return nil, fmt.Errorf("environment variable %v is required", OWMKeyEnvKey)
	}

	logLevel := os.Getenv(LogLevelEnvKey)
	if logLevel == "" {
		logLevel = DefaultLogLevel
	}

	return &Config{
		OWMKey:     owmKey,
		LogLevel:   logLevel,
		OutputJSON: false,
	}, nil
}
