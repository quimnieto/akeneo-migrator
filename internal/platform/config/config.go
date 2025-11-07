package config

import (
	"fmt"

	kit_config "akeneo-migrator/kit/config/static"

	"github.com/spf13/viper"
)

// Config contains the configuration for source and destination
type Config struct {
	AkeneoSource AkeneoSource `json:"akeneoSource" mapstructure:"akeneoSource"`
	AkeneoDest   AkeneoDest   `json:"akeneoDest" mapstructure:"akeneoDest"`
	Source       Source       `json:"source" mapstructure:"source"`
	Dest         Dest         `json:"dest" mapstructure:"dest"`
}

// AkeneoSource contains the source Akeneo configuration from JSON
type AkeneoSource struct {
	API APIConfig `json:"api" mapstructure:"api"`
}

// AkeneoDest contains the destination Akeneo configuration from JSON
type AkeneoDest struct {
	API APIConfig `json:"api" mapstructure:"api"`
}

// APIConfig contains the API configuration
type APIConfig struct {
	URL         string      `json:"url" mapstructure:"url"`
	Credentials Credentials `json:"credentials" mapstructure:"credentials"`
}

// Credentials contains the access credentials
type Credentials struct {
	ClientID string `json:"clientId" mapstructure:"clientId"`
	Secret   string `json:"secret" mapstructure:"secret"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
}

// Source contains the source Akeneo configuration (for compatibility)
type Source struct {
	Host     string `json:"host" mapstructure:"host"`
	ClientID string `json:"clientId" mapstructure:"clientId"`
	Secret   string `json:"secret" mapstructure:"secret"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
}

// Dest contains the destination Akeneo configuration (for compatibility)
type Dest struct {
	Host     string `json:"host" mapstructure:"host"`
	ClientID string `json:"clientId" mapstructure:"clientId"`
	Secret   string `json:"secret" mapstructure:"secret"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
}

// LoadConfig loads the configuration using Viper
func LoadConfig(configLoader kit_config.ConfigurationLoader) (*Config, error) {
	// Get Viper instance for the context
	v := viper.Get("akeneo-migrator").(viper.Viper)

	config := &Config{}

	// Unmarshal the complete configuration
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("error deserializing configuration: %w", err)
	}

	// Map from JSON structure to compatibility structure
	if config.AkeneoSource.API.URL != "" {
		config.Source = Source{
			Host:     config.AkeneoSource.API.URL,
			ClientID: config.AkeneoSource.API.Credentials.ClientID,
			Secret:   config.AkeneoSource.API.Credentials.Secret,
			Username: config.AkeneoSource.API.Credentials.Username,
			Password: config.AkeneoSource.API.Credentials.Password,
		}
	}

	if config.AkeneoDest.API.URL != "" {
		config.Dest = Dest{
			Host:     config.AkeneoDest.API.URL,
			ClientID: config.AkeneoDest.API.Credentials.ClientID,
			Secret:   config.AkeneoDest.API.Credentials.Secret,
			Username: config.AkeneoDest.API.Credentials.Username,
			Password: config.AkeneoDest.API.Credentials.Password,
		}
	}

	// Validate configuration
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfig(config *Config) error {
	// Validate source configuration
	if config.Source.Host == "" || config.Source.ClientID == "" ||
		config.Source.Secret == "" || config.Source.Username == "" ||
		config.Source.Password == "" {
		return fmt.Errorf("incomplete SOURCE configuration")
	}

	// Validate destination configuration
	if config.Dest.Host == "" || config.Dest.ClientID == "" ||
		config.Dest.Secret == "" || config.Dest.Username == "" ||
		config.Dest.Password == "" {
		return fmt.Errorf("incomplete DEST configuration")
	}

	return nil
}
