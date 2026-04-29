package config

import (
	"testing"
)

func TestKafkaDefaults_Topic(t *testing.T) {
	cfg := &Config{}
	kafkaDefaults(cfg)
	if cfg.Kafka.Topic != "portwatch-events" {
		t.Errorf("expected default topic %q, got %q", "portwatch-events", cfg.Kafka.Topic)
	}
}

func TestKafkaDefaults_DoesNotOverride(t *testing.T) {
	cfg := &Config{}
	cfg.Kafka.Topic = "my-topic"
	kafkaDefaults(cfg)
	if cfg.Kafka.Topic != "my-topic" {
		t.Errorf("default should not override existing topic, got %q", cfg.Kafka.Topic)
	}
}

func TestValidateKafka_Disabled(t *testing.T) {
	cfg := &Config{}
	if err := validateKafka(cfg); err != nil {
		t.Fatalf("expected no error when disabled, got: %v", err)
	}
}

func TestValidateKafka_MissingBrokers(t *testing.T) {
	cfg := &Config{}
	cfg.Kafka.Enabled = true
	cfg.Kafka.Topic = "portwatch-events"
	if err := validateKafka(cfg); err == nil {
		t.Fatal("expected error for missing brokers")
	}
}

func TestValidateKafka_MissingTopic(t *testing.T) {
	cfg := &Config{}
	cfg.Kafka.Enabled = true
	cfg.Kafka.Brokers = []string{"localhost:9092"}
	cfg.Kafka.Topic = ""
	if err := validateKafka(cfg); err == nil {
		t.Fatal("expected error for missing topic")
	}
}

func TestValidateKafka_Valid(t *testing.T) {
	cfg := &Config{}
	cfg.Kafka.Enabled = true
	cfg.Kafka.Brokers = []string{"localhost:9092"}
	cfg.Kafka.Topic = "portwatch-events"
	if err := validateKafka(cfg); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}
