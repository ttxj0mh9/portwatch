package config

import "fmt"

// PrometheusConfig holds settings for the Prometheus metrics exporter.
type PrometheusConfig struct {
	Enabled bool   `yaml:"enabled"`
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Path    string `yaml:"path"`
}

func prometheusDefaults() PrometheusConfig {
	return PrometheusConfig{
		Enabled: false,
		Host:    "127.0.0.1",
		Port:    9090,
		Path:    "/metrics",
	}
}

func validatePrometheus(cfg PrometheusConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.Port < 1 || cfg.Port > 65535 {
		return fmt.Errorf("prometheus: port %d is out of range (1-65535)", cfg.Port)
	}
	if cfg.Host == "" {
		return fmt.Errorf("prometheus: host must not be empty")
	}
	if cfg.Path == "" {
		return fmt.Errorf("prometheus: path must not be empty")
	}
	return nil
}
