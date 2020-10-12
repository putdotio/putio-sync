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

// TODO Watch fs events for uploading files larger than 4M

var (
	versionFlag     = flag.Bool("version", false, "print program version")
	printConfigPath = flag.Bool("print-config-path", false, "print config file path")
	configFlag      = flag.String("config", "", "config file path")
)

var config putiosync.Config

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

	var configPath string
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

	err = config.Read(configPath)
	if err != nil && !os.IsNotExist(err) {
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
		os.Exit(exitCodeConfigError) // nolint: gocritic
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
