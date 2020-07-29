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
	versionFlag = flag.Bool("version", false, "print program version")
	debugFlag   = flag.Bool("debug", false, "print debug logs")
	configFlag  = flag.String("config", "", "path of config file")
	username    = flag.String("username", "", "putio account username")
	password    = flag.String("password", "", "putio account password")
)

var (
	configPath     string
	config         Config
	db             *bbolt.DB
	client         *putio.Client
	localPath      string
	remoteFolderID int64
	jobs           []Job
)

func main() {
	flag.Parse()
	if *versionFlag {
		fmt.Println(Version)
		return
	}
	if *debugFlag {
		log.SetLevel(log.DEBUG)
	}
	if *configFlag != "" {
		configPath = *configFlag
	} else {
		var err error
		configPath, err = xdg.ConfigFile(filepath.Join("putio-sync", "config.toml"))
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Infoln("Using config file:", configPath)
	if err := config.Read(); err != nil {
		log.Fatal(err)
	}
	if *username != "" {
		config.Username = *username
	}
	if *password != "" {
		config.Password = *password
	}
	if err := config.Validate(); err != nil {
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
	log.Infoln("Sync finished")
}
