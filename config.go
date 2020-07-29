package main

import (
	"errors"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Username string
	Password string
}

func (c *Config) Read() error {
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		if *configFlag != "" {
			// user specified config must exists
			// default config may not exist, that's okay
			return err
		}
	}
	return nil
}

func (c *Config) Validate() error {
	if c.Username == "" {
		return errors.New("empty username in config")
	}
	if c.Password == "" {
		return errors.New("empty password in config")
	}
	return nil
}
