package config

import "fmt"

// HTTPHandlerConfig holds configuration for the generic HTTP alert handler.
type HTTPHandlerConfig struct {
	Enabled bool              `yaml:"enabled"`
	URL     string            `yaml:"url"`
	Method  string            `yaml:"method"`
	Headers map[string]string `yaml:"headers"`
}

func httpDefaults(cfg *Config) {
	if cfg.Alerts.HTTP.Method == "" {
		cfg.Alerts.HTTP.Method = "POST"
	}
}

func validateHTTP(cfg *Config) error {
	h := cfg.Alerts.HTTP
	if !h.Enabled {
		return nil
	}
	if h.URL == "" {
		return fmt.Errorf("alerts.http.url is required when http alerting is enabled")
	}
	return nil
}
