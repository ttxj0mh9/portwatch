package config

import "fmt"

// MatrixConfig holds settings for the Matrix alert handler.
type MatrixConfig struct {
	Enabled     bool   `yaml:"enabled"`
	Homeserver  string `yaml:"homeserver"`
	AccessToken string `yaml:"access_token"`
	RoomID      string `yaml:"room_id"`
}

// matrixDefaults applies default values to MatrixConfig fields that are unset.
func matrixDefaults(cfg *Config) {
	// No default homeserver, token, or room — all must be user-supplied.
	// This function exists for consistency with other handler config files.
}

// validateMatrix returns an error if the Matrix configuration is enabled but
// incomplete.
func validateMatrix(cfg *Config) error {
	m := cfg.Alerts.Matrix
	if !m.Enabled {
		return nil
	}
	if m.Homeserver == "" {
		return fmt.Errorf("alerts.matrix.homeserver is required when matrix alerts are enabled")
	}
	if m.AccessToken == "" {
		return fmt.Errorf("alerts.matrix.access_token is required when matrix alerts are enabled")
	}
	if m.RoomID == "" {
		return fmt.Errorf("alerts.matrix.room_id is required when matrix alerts are enabled")
	}
	return nil
}
