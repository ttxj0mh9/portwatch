package config

import "fmt"

type RocketChatConfig struct {
	Enabled    bool   `yaml:"enabled"`
	WebhookURL string `yaml:"webhook_url"`
	Channel    string `yaml:"channel"`
	Username   string `yaml:"username"`
}

func rocketChatDefaults(cfg *Config) {
	if cfg.Alerts.RocketChat.Username == "" {
		cfg.Alerts.RocketChat.Username = "portwatch"
	}
}

func validateRocketChat(cfg *Config) error {
	rc := cfg.Alerts.RocketChat
	if !rc.Enabled {
		return nil
	}
	if rc.WebhookURL == "" {
		return fmt.Errorf("alerts.rocketchat.webhook_url is required when rocketchat is enabled")
	}
	return nil
}
