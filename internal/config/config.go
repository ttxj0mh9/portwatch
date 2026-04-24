package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the full portwatch runtime configuration.
type Config struct {
	ScanInterval time.Duration `yaml:"-"`
	IntervalRaw  string        `yaml:"interval"`
	Ports        PortRange     `yaml:"ports"`
	Alert        AlertConfig   `yaml:"alert"`
	Snapshot     SnapshotCfg   `yaml:"snapshot"`
}

type PortRange struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`
}

type AlertConfig struct {
	Log     LogCfg     `yaml:"log"`
	Webhook WebhookCfg `yaml:"webhook"`
	Slack   SlackCfg   `yaml:"slack"`
}

type LogCfg struct {
	File string `yaml:"file"`
}

type WebhookCfg struct {
	URL string `yaml:"url"`
}

type SlackCfg struct {
	WebhookURL string `yaml:"webhook_url"`
}

type SnapshotCfg struct {
	Path    string `yaml:"path"`
	Backups int    `yaml:"backups"`
}

// Default returns a Config populated with sensible defaults.
func Default() Config {
	return Config{
		IntervalRaw:  "30s",
		ScanInterval: 30 * time.Second,
		Ports:        PortRange{Min: 1, Max: 65535},
		Snapshot:     SnapshotCfg{Path: "portwatch.snap", Backups: 3},
	}
}

// Load reads a YAML config file from path and merges it with defaults.
func Load(path string) (Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("config: read file: %w", err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("config: parse yaml: %w", err)
	}

	if cfg.IntervalRaw != "" {
		d, err := time.ParseDuration(cfg.IntervalRaw)
		if err != nil {
			return cfg, fmt.Errorf("config: invalid interval %q: %w", cfg.IntervalRaw, err)
		}
		cfg.ScanInterval = d
	}

	if err := validate(cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func validate(c Config) error {
	if c.ScanInterval <= 0 {
		return errors.New("config: interval must be positive")
	}
	if c.Ports.Min < 1 || c.Ports.Max > 65535 || c.Ports.Min > c.Ports.Max {
		return fmt.Errorf("config: invalid port range %d-%d", c.Ports.Min, c.Ports.Max)
	}
	return nil
}
