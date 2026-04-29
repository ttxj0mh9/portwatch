package config

import "fmt"

// KafkaConfig holds settings for the Kafka alert handler.
type KafkaConfig struct {
	Enabled bool     `yaml:"enabled"`
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
}

func kafkaDefaults(cfg *Config) {
	if cfg.Kafka.Topic == "" {
		cfg.Kafka.Topic = "portwatch-events"
	}
}

func validateKafka(cfg *Config) error {
	if !cfg.Kafka.Enabled {
		return nil
	}
	if len(cfg.Kafka.Brokers) == 0 {
		return fmt.Errorf("kafka: brokers list must not be empty when enabled")
	}
	if cfg.Kafka.Topic == "" {
		return fmt.Errorf("kafka: topic must not be empty when enabled")
	}
	return nil
}
