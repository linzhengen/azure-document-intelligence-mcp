package config

import "github.com/kelseyhightower/envconfig"

// Config holds the application configuration.
type Config struct {
	AzureEndpoint string `envconfig:"AZURE_DOCUMENT_INTELLIGENCE_ENDPOINT" required:"true"`
	AzureAPIKey   string `envconfig:"AZURE_DOCUMENT_INTELLIGENCE_API_KEY" required:"true"`
	Port          int    `envconfig:"PORT" default:"8081"`
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
