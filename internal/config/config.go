package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents the application configuration.
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

// ServerConfig represents the server configuration.
type ServerConfig struct {
	Port         int `yaml:"port"`
	WriteTimeout int `yaml:"writeTimeout"`
	ReadTimeout  int `yaml:"readTimeout"`
}

// DatabaseConfig represents the database configuration.
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

// LoadConfig loads the configuration from a YAML file.
func LoadConfig(filepath string) (*Config, error) {
	// Create a new Config instance with default values
	cfg := &Config{
		Server: ServerConfig{
			Port:         8000,
			WriteTimeout: 15,
			ReadTimeout:  15,
		},
		// Set other default values as necessary
	}

	// Check if the config file exists
	if _, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			// Config file does not exist, return config with default values
			return cfg, nil
		}
		// Error occurred while checking the file
		return nil, err
	}

	// Read the config file
	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the YAML data into the Config struct
	err = yaml.Unmarshal(file, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
