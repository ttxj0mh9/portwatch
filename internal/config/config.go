package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all portwatch configuration.
type Config struct {
	Interval  time.Duration `yaml:"interval"`
	Ports     PortRange     `yaml:"ports"`
	Snapshot  Snapshot      `yaml:"snapshot"`
	Alerts    Alerts        `yaml:"alerts"`
}

// PortRange defines the TCP port scan range.
type PortRange struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`
}

// Snapshot holds snapshot storage settings.
type Snapshot struct {
	Path    string `yaml:"path"`
	Backups int    `yaml:"backups"`
}

// Alerts groups all alerting handler configurations.
type Alerts struct {
	Log       LogConfig        `yaml:"log"`
	File      FileConfig       `yaml:"file"`
	Webhook   WebhookConfig    `yaml:"webhook"`
	Slack     SlackConfig      `yaml:"slack"`
	Email     EmailConfig      `yaml:"email"`
	PagerDuty PagerDutyConfig  `yaml:"pagerduty"`
	OpsGenie  OpsGenieConfig   `yaml:"opsgenie"`
}

type LogConfig struct {
	Enabled bool `yaml:"enabled"`
}

type FileConfig struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

type WebhookConfig struct {
	Enabled bool   `yaml:"enabled"`
	URL     string `yaml:"url"`
}

type SlackConfig struct {
	Enabled    bool   `yaml:"enabled"`
	WebhookURL string `yaml:"webhook_url"`
}

type EmailConfig struct {
	Enabled    bool     `yaml:"enabled"`
	Host       string   `yaml:"host"`
	Port       int      `yaml:"port"`
	From       string   `yaml:"from"`
	Recipients []string `yaml:"recipients"`
}

type PagerDutyConfig struct {
	Enabled    bool   `yaml:"enabled"`
	RoutingKey string `yaml:"routing_key"`
}

// OpsGenieConfig holds OpsGenie alerting configuration.
type OpsGenieConfig struct {
	Enabled bool   `yaml:"enabled"`
	APIKey  string `yaml:"api_key"`
}

// Default returns a Config populated with sensible defaults.
func Default() Config {
	return Config{
		Interval: 30 * time.Second,
		Ports:    PortRange{Min: 1, Max: 65535},
		Snapshot: Snapshot{Path: "/tmp/portwatch_snapshot.json", Backups: 3},
		Alerts:   Alerts{Log: LogConfig{Enabled: true}},
	}
}

// Load reads a YAML config file from path and merges it over defaults.
func Load(path string) (Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("config: read file: %w", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("config: parse yaml: %w", err)
	}

	if err := validate(cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func validate(cfg Config) error {
	if cfg.Interval < time.Second {
		return fmt.Errorf("config: interval must be at least 1s, got %s", cfg.Interval)
	}
	if cfg.Ports.Min < 1 || cfg.Ports.Max > 65535 || cfg.Ports.Min > cfg.Ports.Max {
		return fmt.Errorf("config: invalid port range %d-%d", cfg.Ports.Min, cfg.Ports.Max)
	}
	if cfg.Alerts.OpsGenie.Enabled && cfg.Alerts.OpsGenie.APIKey == "" {
		return fmt.Errorf("config: opsgenie enabled but api_key is empty")
	}
	return nil
}
