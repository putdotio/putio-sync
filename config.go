package putiosync

import (
	"time"
)

type ConfigError struct {
	reason string
}

func newConfigError(reason string) *ConfigError {
	return &ConfigError{reason: reason}
}

func (e *ConfigError) Error() string {
	return "error in config: " + e.reason
}

type Config struct {
	// Username of put.io account.
	Username string
	// Password of put.io account.
	Password string
	// Do not make changes on filesystems.
	DryRun bool
	// Sync repeatedly. Pause given duration between syncs.
	Repeat time.Duration
	// Listen address for HTTP server.
	Server string
}

func (c *Config) validate() error {
	if c.Username == "" {
		return newConfigError("empty username")
	}
	if c.Password == "" {
		return newConfigError("empty password")
	}
	return nil
}
