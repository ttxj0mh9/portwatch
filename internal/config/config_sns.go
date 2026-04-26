package config

import "fmt"

// SNSConfig holds AWS SNS alert settings.
type SNSConfig struct {
	Enabled  bool   `yaml:"enabled"`
	TopicARN string `yaml:"topic_arn"`
	// Region overrides the AWS_REGION env var when set.
	Region string `yaml:"region"`
}

func snsDefaults(cfg *Config) {
	// SNS is opt-in; no defaults to apply beyond zero values.
	_ = cfg
}

func validateSNS(cfg *Config) error {
	if !cfg.Alerts.SNS.Enabled {
		return nil
	}
	if cfg.Alerts.SNS.TopicARN == "" {
		return fmt.Errorf("alerts.sns.topic_arn is required when sns is enabled")
	}
	return nil
}
