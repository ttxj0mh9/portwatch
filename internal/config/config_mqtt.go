package config

import "fmt"

// MQTTConfig holds settings for the MQTT alert handler.
type MQTTConfig struct {
	Enabled   bool   `yaml:"enabled"`
	BrokerURL string `yaml:"broker_url"`
	Topic     string `yaml:"topic"`
	ClientID  string `yaml:"client_id"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
}

func mqttDefaults(cfg *Config) {
	if cfg.Alerts.MQTT.ClientID == "" {
		cfg.Alerts.MQTT.ClientID = "portwatch"
	}
	if cfg.Alerts.MQTT.Topic == "" {
		cfg.Alerts.MQTT.Topic = "portwatch/events"
	}
}

func validateMQTT(cfg *Config) error {
	m := cfg.Alerts.MQTT
	if !m.Enabled {
		return nil
	}
	if m.BrokerURL == "" {
		return fmt.Errorf("alerts.mqtt.broker_url is required when mqtt is enabled")
	}
	return nil
}
