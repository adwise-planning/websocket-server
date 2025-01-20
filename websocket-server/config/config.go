package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	ServerPort  string `mapstructure:"SERVER_PORT"`
}

// LoadConfig reads the configuration from a file or environment variables
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("env")    // required if the config file is not a .json or .yaml file
	viper.AddConfigPath(".")      // optionally look for config in the working directory

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
		return nil, err
	}

	return &config, nil
}