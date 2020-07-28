package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/cenkalti/log"
	"github.com/putdotio/go-putio"
	"go.etcd.io/bbolt"
)

const Version = "0.0.1"

var (
	configFlag  = flag.Bool("config", false, "print config file path")
	versionFlag = flag.Bool("version", false, "print program version")
)

var (
	config         Config
	db             *bbolt.DB
	client         *putio.Client
	localPath      string
	remoteFolderID int64
)

func main() {
	flag.Parse()
	if *versionFlag {
		fmt.Println(Version)
		return
	}
	configPath, err := xdg.ConfigFile(filepath.Join("putio-sync", "config.toml"))
	if err != nil {
		log.Fatal(err)
	}
	if *configFlag {
		fmt.Println(configPath)
		return
	}
	log.Infoln("Using config file:", configPath)
	err = config.Read(configPath)
	if err != nil {
		log.Fatal(err)
	}
	dbPath, err := xdg.DataFile(filepath.Join("putio-sync", "sync.db"))
	if err != nil {
		log.Fatal(err)
	}
	log.Infoln("Using database file:", dbPath)
	db, err = bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = ensureValidClient()
	if err != nil {
		log.Fatal(err)
	}
	err = ensureFolders()
	if err != nil {
		log.Fatal(err)
	}
	err = sync()
	if err != nil {
		log.Fatal(err)
	}
}
