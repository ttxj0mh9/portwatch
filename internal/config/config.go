package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the full portwatch configuration.
type Config struct {
	Ports    []int         `yaml:"ports"`
	Interval time.Duration `yaml:"interval"`
	LogFile  string        `yaml:"log_file"`
	Snapshot string        `yaml:"snapshot"`
	Webhook  WebhookConfig `yaml:"webhook"`
}

// WebhookConfig holds optional webhook alert settings.
type WebhookConfig struct {
	URL     string        `yaml:"url"`
	Timeout time.Duration `yaml:"timeout"`
}

// Default returns a Config populated with sensible defaults.
func Default() Config {
	return Config{
		Ports:    []int{},
		Interval: 30 * time.Second,
		LogFile:  "portwatch.log",
		Snapshot: "portwatch.snap",
		Webhook: WebhookConfig{
			Timeout: 5 * time.Second,
		},
	}
}

// Load reads a YAML config file and merges it with defaults.
func Load(path string) (Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("config: read %s: %w", path, err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("config: parse %s: %w", path, err)
	}

	if err := validate(cfg); err != nil {
		return cfg, fmt.Errorf("config: %w", err)
	}

	return cfg, nil
}

func validate(cfg Config) error {
	if cfg.Interval < time.Second {
		return fmt.Errorf("interval must be at least 1s, got %s", cfg.Interval)
	}
	for _, p := range cfg.Ports {
		if p < 1 || p > 65535 {
			return fmt.Errorf("invalid port %d: must be 1–65535", p)
		}
	}
	return nil
}
