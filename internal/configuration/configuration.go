// Package configuration handles the loading and storing of configuration settings for the application.
// It utilizes the viper library to load settings from environment variables, providing a flexible and
// powerful way to manage configuration in various environments.
package configuration

import (
	"github.com/spf13/viper"
	"time"
)

// Config is a global variable that holds the application's configuration settings.
// It is initialized by the init function which loads the configuration from the environment.
var Config *AppConfig

// init initializes the configuration by loading settings from the environment.
// It calls loadFromEnv to populate the Config variable with the application settings.
func init() {
	Config = loadFromEnv()
}

// AppConfig defines the structure of the application's configuration settings.
// It includes various configurations such as server details, database settings, and cache configurations.
type AppConfig struct {
	Name           string         // Name of the application
	Env            string         // Environment (e.g., development, production)
	HostAddress    string         // Server host address
	HostPort       int            // Server port number
	DocsAddress    string         // Address for API documentation
	DatabaseConfig DatabaseConfig // Configuration settings for the database
	CacheConfig    CacheConfig    // Configuration settings for caching
	Elasticsearch  SearchEngineConfig
}

// DatabaseConfig defines the configuration settings for the database connection.
// It includes details like the database name, driver, DSN, and connection pool settings.
type DatabaseConfig struct {
	Name                   string        // Database name
	Driver                 string        // Database driver type (e.g., mysql, postgres)
	DSN                    string        // Data source name for the database connection
	MaxConnections         int           // Maximum number of open connections to the database
	MaxIdleConnections     int           // Maximum number of idle connections to the database
	MaxLifetimeConnections time.Duration // Maximum amount of time a connection may be reused
}

// CacheConfig defines the configuration for cache connections.
// It includes the data source name for the cache (e.g., Redis).
type CacheConfig struct {
	DSN string // Data source name for the cache
}

type SearchEngineConfig struct {
	URL      string
	User     string
	Password string
}

// loadFromEnv loads configuration settings from environment variables and returns an AppConfig instance.
// It uses viper to handle the environment variables and sets default values if specific configurations are not provided.
func loadFromEnv() *AppConfig {
	viper.AutomaticEnv() // Automatically read environment variables

	return &AppConfig{
		Name:        viper.GetString("APP_NAME"),              // Application name
		Env:         viper.GetString("APP_ENV"),               // Application environment
		HostAddress: viper.GetString("APP_HOST_ADDRESS"),      // Server host address
		HostPort:    viper.GetInt("APP_HOST_PORT"),            // Server port number
		DocsAddress: viper.GetString("APP_DOCS_HOST_ADDRESS"), // Address for API documentation
		DatabaseConfig: DatabaseConfig{
			Driver:                 viper.GetString("MYSQL_DRIVER"),                    // Database driver
			Name:                   viper.GetString("MYSQL_NAME"),                      // Database name
			DSN:                    viper.GetString("MYSQL_DSN"),                       // Data source name
			MaxConnections:         viper.GetInt("MYSQL_MAX_CONNECTIONS"),              // Max open connections
			MaxIdleConnections:     viper.GetInt("MYSQL_MAX_IDLE_CONNECTIONS"),         // Max idle connections
			MaxLifetimeConnections: viper.GetDuration("MYSQL_MAX_CONNECTION_LIFETIME"), // Connection lifetime
		},
		CacheConfig: CacheConfig{
			DSN: viper.GetString("REDIS_DSN"), // Data source name for Redis
		},
		Elasticsearch: SearchEngineConfig{
			URL:      viper.GetString("ES_URL"),
			User:     viper.GetString("ES_USER"),
			Password: viper.GetString("ES_PASSWORD"),
		},
	}
}
