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
	Interval    time.Duration `yaml:"interval"`
	Ports       []int         `yaml:"ports"`
	SnapshotDir string        `yaml:"snapshot_dir"`
	LogFile     string        `yaml:"log_file"`
	Alerts      AlertsConfig  `yaml:"alerts"`
}

// AlertsConfig groups all optional alert handler configurations.
type AlertsConfig struct {
	Webhook   *WebhookConfig   `yaml:"webhook,omitempty"`
	Slack     *SlackConfig     `yaml:"slack,omitempty"`
	Email     *EmailConfig     `yaml:"email,omitempty"`
	PagerDuty *PagerDutyConfig `yaml:"pagerduty,omitempty"`
	OpsGenie  *OpsGenieConfig  `yaml:"opsgenie,omitempty"`
	Discord   *DiscordConfig   `yaml:"discord,omitempty"`
}

type WebhookConfig struct {
	URL string `yaml:"url"`
}

type SlackConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}

type EmailConfig struct {
	Host       string   `yaml:"host"`
	Port       int      `yaml:"port"`
	From       string   `yaml:"from"`
	Recipients []string `yaml:"recipients"`
}

type PagerDutyConfig struct {
	IntegrationKey string `yaml:"integration_key"`
}

type OpsGenieConfig struct {
	APIKey string `yaml:"api_key"`
}

type DiscordConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}

// Default returns a Config populated with sensible defaults.
func Default() Config {
	return Config{
		Interval:    30 * time.Second,
		Ports:       []int{},
		SnapshotDir: "/tmp/portwatch",
		LogFile:     "",
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

	if err := validate(cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func validate(cfg Config) error {
	if cfg.Interval < time.Second {
		return errors.New("config: interval must be at least 1s")
	}
	for _, p := range cfg.Ports {
		if p < 1 || p > 65535 {
			return fmt.Errorf("config: invalid port %d (must be 1-65535)", p)
		}
	}
	if cfg.Alerts.Discord != nil && cfg.Alerts.Discord.WebhookURL == "" {
		return errors.New("config: discord webhook_url must not be empty")
	}
	return nil
}
