package config

// DatadogConfig holds settings for the Datadog alert handler.
type DatadogConfig struct {
	// APIKey is the Datadog API key used for authentication.
	APIKey string `yaml:"api_key"`
}

// datadogDefaults returns a DatadogConfig with zero values (no defaults needed).
func datadogDefaults() DatadogConfig {
	return DatadogConfig{}
}

// validateDatadog returns an error if the DatadogConfig is enabled but invalid.
func validateDatadog(cfg DatadogConfig) error {
	if cfg.APIKey == "" {
		return &ValidationError{Field: "alerts.datadog.api_key", Reason: "api key is required when datadog handler is enabled"}
	}
	return nil
}

// ValidationError represents a configuration validation failure.
type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return "config: invalid field " + e.Field + ": " + e.Reason
}
