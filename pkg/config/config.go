package config

import (
	"github.com/darki73/pac-manager/pkg/config/database"
	"github.com/darki73/pac-manager/pkg/config/pac"
	"github.com/darki73/pac-manager/pkg/config/server"
	"github.com/spf13/viper"
)

// Config represents the configuration of the application
type Config struct {
	// Database represents the database configuration.
	Database *database.Database `mapstructure:"database" yaml:"database" json:"database" toml:"database"`
	// Pac represents the pac configuration.
	Pac *pac.Pac `mapstructure:"pac" yaml:"pac" json:"pac" toml:"pac"`
	// Server represents the server configuration.
	Server *server.Server `mapstructure:"server" yaml:"server" json:"server" toml:"server"`
}

// Initialize returns the configuration of the application.
func Initialize() (*Config, error) {
	config := &Config{
		Database: database.InitializeDefaults(),
		Pac:      pac.InitializeDefaults(),
	}

	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	viper.WatchConfig()

	return config, nil
}

// GetDatabase returns the database configuration.
func (config *Config) GetDatabase() *database.Database {
	return config.Database
}

// GetPac returns the pac configuration.
func (config *Config) GetPac() *pac.Pac {
	return config.Pac
}

// GetServer returns the server configuration.
func (config *Config) GetServer() *server.Server {
	return config.Server
}
