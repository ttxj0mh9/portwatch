package config

import (
	"testing"
)

func TestTelegramDefaults_NoOp(t *testing.T) {
	// telegramDefaults should not overwrite existing values
	cfg := &Config{}
	cfg.Telegram.Token = "existing-token"
	cfg.Telegram.ChatID = "existing-chat-id"

	telegramDefaults(cfg)

	if cfg.Telegram.Token != "existing-token" {
		t.Errorf("expected token to remain 'existing-token', got %q", cfg.Telegram.Token)
	}
	if cfg.Telegram.ChatID != "existing-chat-id" {
		t.Errorf("expected chat_id to remain 'existing-chat-id', got %q", cfg.Telegram.ChatID)
	}
}

func TestValidateTelegram_Disabled(t *testing.T) {
	cfg := &Config{}
	cfg.Telegram.Enabled = false

	if err := validateTelegram(cfg); err != nil {
		t.Errorf("expected no error when telegram disabled, got: %v", err)
	}
}

func TestValidateTelegram_MissingToken(t *testing.T) {
	cfg := &Config{}
	cfg.Telegram.Enabled = true
	cfg.Telegram.Token = ""
	cfg.Telegram.ChatID = "123456789"

	if err := validateTelegram(cfg); err == nil {
		t.Error("expected error for missing token, got nil")
	}
}

func TestValidateTelegram_MissingChatID(t *testing.T) {
	cfg := &Config{}
	cfg.Telegram.Enabled = true
	cfg.Telegram.Token = "bot123:ABC"
	cfg.Telegram.ChatID = ""

	if err := validateTelegram(cfg); err == nil {
		t.Error("expected error for missing chat_id, got nil")
	}
}

func TestValidateTelegram_Valid(t *testing.T) {
	cfg := &Config{}
	cfg.Telegram.Enabled = true
	cfg.Telegram.Token = "bot123:ABC"
	cfg.Telegram.ChatID = "123456789"

	if err := validateTelegram(cfg); err != nil {
		t.Errorf("expected no error for valid telegram config, got: %v", err)
	}
}

func TestValidateTelegram_EnabledWithBothFields(t *testing.T) {
	cfg := &Config{}
	cfg.Telegram.Enabled = true
	cfg.Telegram.Token = "bot987:XYZ"
	cfg.Telegram.ChatID = "-100123456789"

	if err := validateTelegram(cfg); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}
