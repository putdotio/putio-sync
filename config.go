package putiosync

import (
	"strings"
	"time"
)

type ConfigError struct {
	Reason string
}

func newConfigError(reason string) *ConfigError {
	return &ConfigError{Reason: reason}
}

func (e *ConfigError) Error() string {
	return "error in config: " + e.Reason
}

type Config struct {
	// Username of put.io account.
	Username string
	// Password of put.io account.
	// An OAuth token can be used instead of password.
	// Token must be prefixed with "token/" (without quotes).
	// In that case, Username is not required.
	Password string
	// Do not make changes on filesystems. Only calculate what needs to be done.
	DryRun bool
	// Sync repeatedly. Pause given duration between syncs.
	// If value is greater than zero, Sync function does not return on synchronization errors.
	// However, it retruns on authentication error since there is no point in retrying in this case.
	Repeat time.Duration
	// Listen address for HTTP server.
	// The server has an endpoint for getting the status of the sync operation.
	Server string
	// Set log level to debug.
	Debug bool
}

func (c *Config) validate() error {
	if c.Username == "" && !strings.HasPrefix(c.Password, "token/") {
		return newConfigError("empty username")
	}
	if c.Password == "" {
		return newConfigError("empty password")
	}
	return nil
}
