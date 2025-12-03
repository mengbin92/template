// Package config provides configuration loading functionality.
// It uses Viper to load configuration from YAML files.
package config

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// LoadConfig loads the application configuration from the config file.
// The default config file path is "./config/config.yaml".
//
// The function will exit the program if:
//   - The config file cannot be found
//   - The config file cannot be read
//   - The config file format is invalid
func LoadConfig() {
	viper.SetConfigFile("./config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		errMsg := fmt.Sprintf("load config error: %s", err.Error())
		fmt.Fprintf(os.Stderr, "%s\n", errMsg)
		os.Exit(1)
	}
}

// LoadConfigWithPath loads the application configuration from a specified config file path.
//
// Parameters:
//   - configPath: The path to the configuration file
//
// Returns:
//   - error: Error if the config file cannot be loaded
func LoadConfigWithPath(configPath string) error {
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, "load config error")
	}
	return nil
}