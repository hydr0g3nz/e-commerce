package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds all the configurations from the YAML file.
type Config struct {
	AppName  string         `mapstructure:"app_name"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"db"`
	Upload   UploadConfig   `mapstructure:"upload"`
	Key      KeyConfig      `mapstructure:"key"`
}

// ServerConfig holds server-related configurations.
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Path string `mapstructure:"path"`
}
type KeyConfig struct {
	AccessToken  string `mapstructure:"access_token"`
	RefreshToken string `mapstructure:"refresh_token"`
}

// DatabaseConfig holds database-related configurations.
type DatabaseConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Name     string `mapstructure:"name"`
	Port     string `mapstructure:"port"`
}
type UploadConfig struct {
	UploadPath string `mapstructure:"upload_path"`
	ServerPath string `mapstructure:"server_path"`
}

// LoadConfig loads the configuration from the YAML file and unmarshals it into the Config struct.
func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path) // Set the file path, e.g., "./config.yml"
	viper.SetConfigType("yml")
	// viper.AutomaticEnv() // Automatically use environment variables where applicable

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}
	fmt.Println("config", config)
	return &config, nil
}
