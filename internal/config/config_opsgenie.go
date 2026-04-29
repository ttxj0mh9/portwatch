package config

import "fmt"

type OpsGenieConfig struct {
	Enabled bool     `yaml:"enabled"`
	APIKey  string   `yaml:"api_key"`
	APIURL  string   `yaml:"api_url"`
	Tags    []string `yaml:"tags"`
}

func opsgenieDefaults(cfg *Config) {
	if cfg.Alerts.OpsGenie.APIURL == "" {
		cfg.Alerts.OpsGenie.APIURL = "https://api.opsgenie.com/v2/alerts"
	}
}

func validateOpsGenie(cfg *Config) error {
	og := cfg.Alerts.OpsGenie
	if !og.Enabled {
		return nil
	}
	if og.APIKey == "" {
		return fmt.Errorf("alerts.opsgenie.api_key is required when opsgenie is enabled")
	}
	return nil
}
