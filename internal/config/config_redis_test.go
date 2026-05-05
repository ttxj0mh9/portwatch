package config

import (
	"testing"
)

func TestRedisDefaults_Addr(t *testing.T) {
	cfg := &Config{}
	redisDefaults(cfg)
	if cfg.Alerts.Redis.Addr != "localhost:6379" {
		t.Errorf("addr = %q, want \"localhost:6379\"", cfg.Alerts.Redis.Addr)
	}
}

func TestRedisDefaults_Channel(t *testing.T) {
	cfg := &Config{}
	redisDefaults(cfg)
	if cfg.Alerts.Redis.Channel != "portwatch:events" {
		t.Errorf("channel = %q, want \"portwatch:events\"", cfg.Alerts.Redis.Channel)
	}
}

func TestRedisDefaults_DoesNotOverride(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.Redis.Addr = "redis.internal:6380"
	cfg.Alerts.Redis.Channel = "custom:channel"
	redisDefaults(cfg)
	if cfg.Alerts.Redis.Addr != "redis.internal:6380" {
		t.Errorf("addr should not be overridden")
	}
	if cfg.Alerts.Redis.Channel != "custom:channel" {
		t.Errorf("channel should not be overridden")
	}
}

func TestValidateRedis_Disabled(t *testing.T) {
	cfg := &Config{}
	if err := validateRedis(cfg); err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestValidateRedis_MissingAddr(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.Redis.Enabled = true
	cfg.Alerts.Redis.Channel = "portwatch:events"
	// Addr intentionally empty
	if err := validateRedis(cfg); err == nil {
		t.Error("expected error for missing addr")
	}
}

func TestValidateRedis_MissingChannel(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.Redis.Enabled = true
	cfg.Alerts.Redis.Addr = "localhost:6379"
	if err := validateRedis(cfg); err == nil {
		t.Error("expected error for missing channel")
	}
}

func TestValidateRedis_Valid(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.Redis.Enabled = true
	cfg.Alerts.Redis.Addr = "localhost:6379"
	cfg.Alerts.Redis.Channel = "portwatch:events"
	if err := validateRedis(cfg); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
