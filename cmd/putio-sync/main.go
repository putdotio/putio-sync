package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/adrg/xdg"
	"github.com/cenkalti/log"
	putiosync "github.com/putdotio/putio-sync/v2"
)

const (
	exitCodeConfigError        = 10
	exitCodeInvalidCredentials = 11
)

// These variables are set by goreleaser on build.
var (
	version = "0.0.0"
	commit  = ""
	date    = ""
)

// TODO Watch fs events for files larger than 4M
// TODO Move Uploader inside go-putio pkg
// TODO Add debug flag to config

var (
	versionFlag = flag.Bool("version", false, "print program version")
	debugFlag   = flag.Bool("debug", false, "enable debug logs")
	configFlag  = flag.String("config", "", "path of config file")
	username    = flag.String("username", "", "put.io account username")
	password    = flag.String("password", "", "put.io account password")
	dryrun      = flag.Bool("dryrun", false, "do not make changes on filesystems")
	repeat      = flag.Duration("repeat", 0, "sync repeatedly, pause given duration between syncs")
	server      = flag.String("server", "", "listen address for HTTP API")

	printConfigPath = flag.Bool("print-config-path", false, "print config path")
)

var (
	configPath string
	config     putiosync.Config
)

func readConfig(configPath string, mustExist bool) error {
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		if mustExist {
			// user specified config must exists
			// default config may not exist, that's okay
			return err
		}
	}
	if *username != "" {
		config.Username = *username
	}
	if *password != "" {
		config.Password = *password
	}
	if *dryrun {
		config.DryRun = *dryrun
	}
	if *repeat != 0 {
		config.Repeat = *repeat
	}
	if *server != "" {
		config.Server = *server
	}
	return nil
}

func versionString() string {
	if len(commit) > 7 {
		commit = commit[:7]
	}
	return fmt.Sprintf("%s (%s) [%s]", version, commit, date)
}

func main() {
	var err error
	flag.Parse()
	if *versionFlag {
		fmt.Println(versionString())
		return
	}
	if *debugFlag {
		log.SetLevel(log.DEBUG)
	}
	if *configFlag != "" {
		configPath = *configFlag
	} else {
		configPath, err = xdg.ConfigFile(filepath.Join("putio-sync", "config.toml"))
		if err != nil {
			log.Fatal(err)
		}
	}
	if *printConfigPath {
		fmt.Println(configPath)
		return
	}
	log.Infoln("Starting putio-sync version", versionString())
	log.Infof("Using config file %q", configPath)
	if err = readConfig(configPath, *configFlag != ""); err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-ch
		log.Noticef("Received %s. Stopping sync.", sig)
		cancel()
	}()
	err = putiosync.Sync(ctx, config)
	var configError *putiosync.ConfigError
	if errors.As(err, &configError) {
		fmt.Fprintln(os.Stderr, configError.Reason)
		os.Exit(exitCodeConfigError)
		return
	}
	if errors.Is(err, putiosync.ErrInvalidCredentials) {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(exitCodeInvalidCredentials)
		return
	}
	if err != nil {
		log.Fatal(err)
	}
}
