package putiosync

import (
	"errors"
	"time"
)

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
		return errors.New("empty username in config")
	}
	if c.Password == "" {
		return errors.New("empty password in config")
	}
	return nil
}
