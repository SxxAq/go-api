// Package config handles loading application configuration from YAML files
// and environment variables using the cleanenv library.
package config

import (
	"flag" // For parsing command-line flags
	"log"  // For logging errors and exiting program
	"os"   // For accessing environment variables and checking file existence

	"github.com/ilyakaznacheev/cleanenv" // Third-party package for config parsing
)

// HttpServer holds HTTP server-specific configuration.
type HttpServer struct {
	Addr string `yaml:"addr"` // Maps to 'addr' key in YAML
}

// Config is the main application configuration struct.
// It can be populated from a YAML file or environment variables.
type Config struct {
	Env         string               `yaml:"env" env:"ENV" env-required:"true"` // Environment (e.g., dev, prod), required
	StoragePath string               `yaml:"storage_path" env-required:"true"`  // Path for storing files, required
	HttpServer  `yaml:"http_server"` // Embedded struct for HTTP server config
}

// MustLoad loads the configuration from environment variable, command-line flag, or YAML file.
// It stops the program immediately if anything goes wrong (fail-fast pattern).
func MustLoad() *Config {
	var cfgPath string

	// 1. Check if CONFIG_PATH environment variable is set
	cfgPath = os.Getenv("CONFIG_PATH")

	// 2. If not set, check for -config command-line flag
	if cfgPath == "" {
		// Define a command-line flag "config"
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse() // Parse all command-line flags

		cfgPath = *flags // Use flag value if provided
		if cfgPath == "" {
			// If neither ENV nor flag is set, stop program
			log.Fatal("Config path is not set. Use CONFIG_PATH env or -config flag")
		}
	}

	// 3. Check if the file exists
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", cfgPath)
	}

	// 4. Initialize an empty Config struct
	var cfg Config

	// 5. Use cleanenv to read YAML file and populate the struct
	//    - Fields can also be overridden by environment variables
	//    - Fields marked with env-required:"true" must have values
	err := cleanenv.ReadConfig(cfgPath, &cfg)
	if err != nil {
		log.Fatalf("Cannot read config file: %s", err.Error())
	}

	// 6. Return pointer to populated Config struct
	return &cfg
}
