package config

import "fmt"

// VictorOpsConfig holds settings for the VictorOps (Splunk On-Call) alert handler.
type VictorOpsConfig struct {
	Enabled    bool   `yaml:"enabled"`
	RESTURL    string `yaml:"rest_url"`
	RoutingKey string `yaml:"routing_key"`
}

func victorOpsDefaults(cfg *Config) {
	// VictorOps is disabled by default; no additional field defaults needed.
	if cfg.Alerts.VictorOps.RoutingKey == "" {
		cfg.Alerts.VictorOps.RoutingKey = "default"
	}
}

func validateVictorOps(cfg *Config) error {
	v := cfg.Alerts.VictorOps
	if !v.Enabled {
		return nil
	}
	if v.RESTURL == "" {
		return fmt.Errorf("alerts.victorops.rest_url is required when victorops is enabled")
	}
	if v.RoutingKey == "" {
		return fmt.Errorf("alerts.victorops.routing_key is required when victorops is enabled")
	}
	return nil
}
