package config

import "fmt"

// TwilioConfig holds settings for the Twilio SMS alert handler.
type TwilioConfig struct {
	Enabled    bool     `yaml:"enabled"`
	AccountSID string   `yaml:"account_sid"`
	AuthToken  string   `yaml:"auth_token"`
	From       string   `yaml:"from_number"`
	To         []string `yaml:"to_numbers"`
}

// twilioDefaults applies default values to TwilioConfig.
// Currently there are no sensible defaults beyond what is required.
func twilioDefaults(cfg *Config) {
	// No automatic defaults; all fields are user-supplied credentials.
}

// validateTwilio returns an error if the Twilio configuration is invalid.
func validateTwilio(cfg *Config) error {
	t := cfg.Alerts.Twilio
	if !t.Enabled {
		return nil
	}
	if t.AccountSID == "" {
		return fmt.Errorf("alerts.twilio.account_sid is required when twilio is enabled")
	}
	if t.AuthToken == "" {
		return fmt.Errorf("alerts.twilio.auth_token is required when twilio is enabled")
	}
	if t.From == "" {
		return fmt.Errorf("alerts.twilio.from_number is required when twilio is enabled")
	}
	if len(t.To) == 0 {
		return fmt.Errorf("alerts.twilio.to_numbers must contain at least one number when twilio is enabled")
	}
	return nil
}
