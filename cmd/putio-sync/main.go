package main

import (
	"context"
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

// These variables are set by goreleaser on build.
var (
	version = "0.0.0"
	commit  = ""
	date    = ""
)

// TODO Watch fs events for files larger than 4M

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

func main() {
	var err error
	flag.Parse()
	if *versionFlag {
		if len(commit) > 7 {
			commit = commit[:7]
		}
		fmt.Printf("%s (%s) [%s]\n", version, commit, date)
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
	if err = putiosync.Sync(ctx, config); err != nil {
		log.Fatal(err)
	}
}
