package config

import "fmt"

type TwilioConfig struct {
	Enabled    bool   `yaml:"enabled"`
	AccountSID string `yaml:"account_sid"`
	AuthToken  string `yaml:"auth_token"`
	FromNumber string `yaml:"from_number"`
	ToNumber   string `yaml:"to_number"`
}

func twilioDefaults(cfg *Config) {
	if cfg.Alerts.Twilio.FromNumber == "" {
		cfg.Alerts.Twilio.FromNumber = ""
	}
}

func validateTwilio(cfg *Config) error {
	t := cfg.Alerts.Twilio
	if !t.Enabled {
		return nil
	}
	if t.AccountSID == "" {
		return fmt.Errorf("twilio: account_sid is required when enabled")
	}
	if t.AuthToken == "" {
		return fmt.Errorf("twilio: auth_token is required when enabled")
	}
	if t.FromNumber == "" {
		return fmt.Errorf("twilio: from_number is required when enabled")
	}
	if t.ToNumber == "" {
		return fmt.Errorf("twilio: to_number is required when enabled")
	}
	return nil
}
