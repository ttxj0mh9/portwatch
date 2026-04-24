package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the portwatch daemon configuration.
type Config struct {
	ScanInterval time.Duration `yaml:"scan_interval"`
	Ports        []int         `yaml:"ports"`
	AlertOnNew   bool          `yaml:"alert_on_new"`
	AlertOnClosed bool         `yaml:"alert_on_closed"`
	LogFile      string        `yaml:"log_file"`
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		ScanInterval:  30 * time.Second,
		Ports:         []int{},
		AlertOnNew:    true,
		AlertOnClosed: true,
		LogFile:       "",
	}
}

// Load reads and parses a YAML config file from the given path.
// Missing fields fall back to defaults.
func Load(path string) (*Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("config: parse yaml: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that the configuration values are sensible.
func (c *Config) Validate() error {
	if c.ScanInterval < time.Second {
		return fmt.Errorf("config: scan_interval must be at least 1s, got %s", c.ScanInterval)
	}
	for _, p := range c.Ports {
		if p < 1 || p > 65535 {
			return fmt.Errorf("config: port %d is out of valid range (1-65535)", p)
		}
	}
	return nil
}
