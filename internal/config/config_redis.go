package config

import "fmt"

// RedisConfig holds configuration for the Redis pub/sub alert handler.
type RedisConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	Channel  string `yaml:"channel"`
}

func redisDefaults(cfg *Config) {
	if cfg.Alerts.Redis.Channel == "" {
		cfg.Alerts.Redis.Channel = "portwatch:events"
	}
	if cfg.Alerts.Redis.Addr == "" {
		cfg.Alerts.Redis.Addr = "localhost:6379"
	}
}

func validateRedis(cfg *Config) error {
	r := cfg.Alerts.Redis
	if !r.Enabled {
		return nil
	}
	if r.Addr == "" {
		return fmt.Errorf("alerts.redis.addr is required when redis is enabled")
	}
	if r.Channel == "" {
		return fmt.Errorf("alerts.redis.channel is required when redis is enabled")
	}
	return nil
}
