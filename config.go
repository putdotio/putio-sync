package main

import (
	"errors"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Username       string
	Password       string
	RemoteFolderID int64
	LocalPath      string
}

func (c *Config) Read(configPath string) error {
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return err
	}
	if c.Username == "" {
		return errors.New("empty username in config")
	}
	if c.Password == "" {
		return errors.New("empty password in config")
	}
	return nil
}
