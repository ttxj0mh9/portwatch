package config

import "fmt"

type ZulipConfig struct {
	Enabled    bool   `yaml:"enabled"`
	BaseURL    string `yaml:"base_url"`
	Email      string `yaml:"email"`
	APIKey     string `yaml:"api_key"`
	Stream     string `yaml:"stream"`
	Topic      string `yaml:"topic"`
}

func zulipDefaults(cfg *Config) {
	if cfg.Alerts.Zulip.Topic == "" {
		cfg.Alerts.Zulip.Topic = "portwatch"
	}
}

func validateZulip(cfg *Config) error {
	z := cfg.Alerts.Zulip
	if !z.Enabled {
		return nil
	}
	if z.BaseURL == "" {
		return fmt.Errorf("alerts.zulip.base_url is required when zulip is enabled")
	}
	if z.Email == "" {
		return fmt.Errorf("alerts.zulip.email is required when zulip is enabled")
	}
	if z.APIKey == "" {
		return fmt.Errorf("alerts.zulip.api_key is required when zulip is enabled")
	}
	if z.Stream == "" {
		return fmt.Errorf("alerts.zulip.stream is required when zulip is enabled")
	}
	return nil
}
