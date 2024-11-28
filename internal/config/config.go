// internal/config/config.go

package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds the application configuration.
type Config struct {
	AppPort string
}

// LoadConfig loads configuration from a specified file path, environment variables, and/or config files.
func LoadConfig(filePath string) (*Config, error) {
	if filePath != "" {
		viper.SetConfigFile(filePath)
	} else {
		viper.SetConfigName("config") // Config file name (without extension)
		viper.AddConfigPath(".")      // Look for the config file in the current directory
	}

	viper.SetConfigType("yaml") // Config file type (can be JSON, TOML, etc.)

	// Optionally, you can also set environment variables prefix
	viper.SetEnvPrefix("MYAPP")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	config := &Config{}

	// Unmarshal the configuration into the Config struct
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %v", err)
	}

	return config, nil
}
