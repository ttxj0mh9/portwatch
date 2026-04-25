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
	ScanInterval time.Duration `yaml:"-"`
	ScanIntervalRaw string    `yaml:"scan_interval"`

	Ports PortsConfig `yaml:"ports"`

	Snapshot SnapshotConfig `yaml:"snapshot"`

	Alerts AlertsConfig `yaml:"alerts"`
}

type PortsConfig struct {
	Include []int `yaml:"include"`
	Exclude []int `yaml:"exclude"`
}

type SnapshotConfig struct {
	Path        string `yaml:"path"`
	BackupCount int    `yaml:"backup_count"`
}

type AlertsConfig struct {
	Log      *LogConfig      `yaml:"log"`
	File     *FileConfig     `yaml:"file"`
	Webhook  *WebhookConfig  `yaml:"webhook"`
	Slack    *SlackConfig    `yaml:"slack"`
	Email    *EmailConfig    `yaml:"email"`
	PagerDuty *PagerDutyConfig `yaml:"pagerduty"`
	OpsGenie *OpsGenieConfig `yaml:"opsgenie"`
	Discord  *DiscordConfig  `yaml:"discord"`
	Teams    *TeamsConfig    `yaml:"teams"`
}

type LogConfig struct{}
type FileConfig struct{ Path string `yaml:"path"` }
type WebhookConfig struct{ URL string `yaml:"url"` }
type SlackConfig struct{ WebhookURL string `yaml:"webhook_url"` }
type EmailConfig struct {
	Host       string   `yaml:"host"`
	Port       int      `yaml:"port"`
	From       string   `yaml:"from"`
	Recipients []string `yaml:"recipients"`
}
type PagerDutyConfig struct{ IntegrationKey string `yaml:"integration_key"` }
type OpsGenieConfig struct{ APIKey string `yaml:"api_key"` }
type DiscordConfig struct{ WebhookURL string `yaml:"webhook_url"` }
type TeamsConfig struct{ WebhookURL string `yaml:"webhook_url"` }

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		ScanInterval:    30 * time.Second,
		ScanIntervalRaw: "30s",
		Snapshot: SnapshotConfig{
			Path:        "/var/lib/portwatch/snapshot.json",
			BackupCount: 3,
		},
		Alerts: AlertsConfig{
			Log: &LogConfig{},
		},
	}
}

// Load reads a YAML config file and merges it with defaults.
func Load(path string) (*Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("config: parse yaml: %w", err)
	}

	if cfg.ScanIntervalRaw != "" {
		d, err := time.ParseDuration(cfg.ScanIntervalRaw)
		if err != nil {
			return nil, fmt.Errorf("config: invalid scan_interval %q: %w", cfg.ScanIntervalRaw, err)
		}
		cfg.ScanInterval = d
	}

	if err := validate(cfg); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	return cfg, nil
}

func validate(cfg *Config) error {
	if cfg.ScanInterval <= 0 {
		return errors.New("scan_interval must be positive")
	}
	for _, p := range cfg.Ports.Include {
		if p < 1 || p > 65535 {
			return fmt.Errorf("invalid port in include list: %d", p)
		}
	}
	for _, p := range cfg.Ports.Exclude {
		if p < 1 || p > 65535 {
			return fmt.Errorf("invalid port in exclude list: %d", p)
		}
	}
	if cfg.Snapshot.BackupCount < 0 {
		return errors.New("snapshot.backup_count must be non-negative")
	}
	return nil
}
