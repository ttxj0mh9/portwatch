package config

import "fmt"

// SplunkConfig holds settings for the Splunk HEC alert handler.
type SplunkConfig struct {
	Enabled bool   `yaml:"enabled"`
	HECURL  string `yaml:"hec_url"`
	Token   string `yaml:"token"`
	Source  string `yaml:"source"`
}

func splunkDefaults(cfg *Config) {
	if cfg.Alerts.Splunk.Source == "" {
		cfg.Alerts.Splunk.Source = "portwatch"
	}
}

func validateSplunk(cfg *Config) error {
	s := cfg.Alerts.Splunk
	if !s.Enabled {
		return nil
	}
	if s.HECURL == "" {
		return fmt.Errorf("alerts.splunk.hec_url is required when splunk is enabled")
	}
	if s.Token == "" {
		return fmt.Errorf("alerts.splunk.token is required when splunk is enabled")
	}
	return nil
}
