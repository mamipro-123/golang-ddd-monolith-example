package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// DatabaseConfig holds the database connection parameters
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`              // Corresponds to the 'host' key in YAML
	Port            int           `mapstructure:"port"`              // Corresponds to the 'port' key in YAML
	User            string        `mapstructure:"user"`              // Corresponds to the 'user' key in YAML
	Password        string        `mapstructure:"password"`          // Corresponds to the 'password' key in YAML
	DBName          string        `mapstructure:"dbname"`            // Corresponds to the 'dbname' key in YAML
	SSLMode         string        `mapstructure:"sslmode"`           // Corresponds to the 'sslmode' key in YAML
	MaxOpenConns    int           `mapstructure:"max_open_conns"`    // Corresponds to the 'max_open_conns' key in YAML
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`    // Corresponds to the 'max_idle_conns' key in YAML
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"` // Corresponds to the 'conn_max_lifetime' key in YAML (e.g. "5m")
}

// SMTPConfig holds SMTP server details
type SMTPConfig struct {
	Host   string `mapstructure:"host"`
	Port   int    `mapstructure:"port"`
	User   string `mapstructure:"user"`
	Pass   string `mapstructure:"pass"`
	From   string `mapstructure:"from"`
	Secure bool   `mapstructure:"secure"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port        string `mapstructure:"port"`
	MetricsPort string `mapstructure:"metrics_port"`
}

// Config holds the general application configuration
type Config struct {
	SMTP     SMTPConfig     `mapstructure:"smtp"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"` // <-- NEWLY ADDED FIELD
}

// GlobalConfig is the global configuration variable
var GlobalConfig Config

// LoadConfig reads the configuration file and loads it into the struct
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config") // config.yaml or config.json etc. will be searched
	viper.SetConfigType("yaml")   // Specifies the file type

	// Find the location of the config file
	configPath := findConfigPath()
	if configPath == "" {
		// If config.yaml is not found in the specified paths, we can also try the project root directory
		// Or we can directly return an error
		// Finding the project root directory can be a bit more complex,
		// so for now, we are only searching specific paths.
		// If necessary, Viper's AddConfigPath function can also be used.
		// viper.AddConfigPath(".") // Working directory
		// viper.AddConfigPath("./config") // ./config directory
		return nil, fmt.Errorf("config.yaml file not found in any of the specified search paths")
	}
	fmt.Printf("Using config file: %s\n", configPath) // Log which config file is being used
	viper.SetConfigFile(configPath)                   // Use the found config file

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		// Specifically check for file not found error
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found at %s: %w", configPath, err)
		}
		return nil, fmt.Errorf("error reading config file '%s': %w", configPath, err)
	}

	// Also read environment variables (optional, useful for prioritization)
	viper.AutomaticEnv()
	// Example: DATABASE_HOST environment variable can override config.database.host
	// viper.SetEnvPrefix("APP") // You can set a prefix for environment variables (e.g. APP_DATABASE_HOST)
	// viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // for config.database.host -> DATABASE_HOST conversion

	// Convert config to struct
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return nil, fmt.Errorf("error unmarshaling config into struct: %w", err)
	}

	return &GlobalConfig, nil
}

// GetConfig returns the loaded global configuration
func GetConfig() *Config {
	return &GlobalConfig
}

// findConfigPath searches for the config.yaml file in various locations
func findConfigPath() string {
	searchPaths := []string{
		"./config", // /config folder inside the project
		".",        // Project root directory / working directory
		// Add other common locations if necessary
		// "/etc/myapp/",
		// filepath.Join(os.Getenv("HOME"), ".myapp"),
	}

	for _, path := range searchPaths {
		fullPath := filepath.Join(path, "config.yaml")
		if _, err := os.Stat(fullPath); err == nil {
			// File found and accessible
			return fullPath
		}
	}

	// Not found anywhere
	return ""
}
