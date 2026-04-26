package config

// TelegramConfig holds configuration for the Telegram alert handler.
type TelegramConfig struct {
	Enabled bool   `yaml:"enabled"`
	Token   string `yaml:"token"`
	ChatID  string `yaml:"chat_id"`
}

// telegramDefaults applies default values to the Telegram config.
// Currently no defaults are needed beyond zero values, but the function
// is kept for consistency with other handler config patterns.
func telegramDefaults(cfg *Config) {
	// No default values to apply for Telegram.
}

// validateTelegram checks that the Telegram configuration is valid.
// If Telegram alerting is disabled, no validation is performed.
func validateTelegram(cfg *Config) error {
	if !cfg.Alerts.Telegram.Enabled {
		return nil
	}
	if cfg.Alerts.Telegram.Token == "" {
		return newValidationError("alerts.telegram.token", "must not be empty when telegram is enabled")
	}
	if cfg.Alerts.Telegram.ChatID == "" {
		return newValidationError("alerts.telegram.chat_id", "must not be empty when telegram is enabled")
	}
	return nil
}
