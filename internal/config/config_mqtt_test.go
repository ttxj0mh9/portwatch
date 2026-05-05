package config

import (
	"testing"
)

func TestMQTTDefaults_ClientID(t *testing.T) {
	cfg := &Config{}
	mqttDefaults(cfg)
	if cfg.Alerts.MQTT.ClientID != "portwatch" {
		t.Errorf("expected default client_id 'portwatch', got %q", cfg.Alerts.MQTT.ClientID)
	}
}

func TestMQTTDefaults_Topic(t *testing.T) {
	cfg := &Config{}
	mqttDefaults(cfg)
	if cfg.Alerts.MQTT.Topic != "portwatch/events" {
		t.Errorf("expected default topic 'portwatch/events', got %q", cfg.Alerts.MQTT.Topic)
	}
}

func TestMQTTDefaults_DoesNotOverride(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.MQTT.ClientID = "custom-id"
	cfg.Alerts.MQTT.Topic = "my/topic"
	mqttDefaults(cfg)
	if cfg.Alerts.MQTT.ClientID != "custom-id" {
		t.Errorf("expected client_id to remain 'custom-id', got %q", cfg.Alerts.MQTT.ClientID)
	}
	if cfg.Alerts.MQTT.Topic != "my/topic" {
		t.Errorf("expected topic to remain 'my/topic', got %q", cfg.Alerts.MQTT.Topic)
	}
}

func TestValidateMQTT_Disabled(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.MQTT.Enabled = false
	if err := validateMQTT(cfg); err != nil {
		t.Errorf("unexpected error for disabled mqtt: %v", err)
	}
}

func TestValidateMQTT_MissingBrokerURL(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.MQTT.Enabled = true
	if err := validateMQTT(cfg); err == nil {
		t.Error("expected error for missing broker_url")
	}
}

func TestValidateMQTT_Valid(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.MQTT.Enabled = true
	cfg.Alerts.MQTT.BrokerURL = "tcp://localhost:1883"
	if err := validateMQTT(cfg); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
