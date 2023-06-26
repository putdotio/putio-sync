package putiosync

import (
	"os"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
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
	// Sync files to/from this dir in computer.
	LocalDir string
	// Do not make changes on filesystems. Only calculate what needs to be done.
	DryRun bool
	// Stop after first sync operation.
	// Otherwise, sync operation is repeated continuously while waiting for some duration between syncs.
	// If there is authentication error, Sync also returns since there is no point in retrying in this case.
	Once bool
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

func (c *Config) Read(configPath string) error {
	k := koanf.New(".")
	err := k.Load(file.Provider(configPath), toml.Parser())
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	err = k.Load(env.Provider("PUTIO_", ".", func(s string) string {
		return strings.ReplaceAll(strings.TrimPrefix(s, "PUTIO_"), "_", ".")
	}), nil)
	if err != nil {
		return err
	}
	err = k.Unmarshal("", c)
	if err != nil {
		return err
	}
	c.setDefaults()
	return nil
}

func (c *Config) setDefaults() {
	if c.LocalDir == "" {
		c.LocalDir = "~/putio-sync"
	}
}
