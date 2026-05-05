package config

import "fmt"

// LokiConfig holds configuration for the Grafana Loki alert handler.
type LokiConfig struct {
	Enabled  bool              `yaml:"enabled"`
	PushURL  string            `yaml:"push_url"`
	Labels   map[string]string `yaml:"labels"`
}

func lokiDefaults(cfg *Config) {
	if cfg.Loki.Labels == nil {
		cfg.Loki.Labels = map[string]string{"app": "portwatch"}
	}
}

func validateLoki(cfg *Config) error {
	if !cfg.Loki.Enabled {
		return nil
	}
	if cfg.Loki.PushURL == "" {
		return fmt.Errorf("loki: push_url is required when loki is enabled")
	}
	return nil
}
