package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/cenkalti/log"
	"github.com/putdotio/go-putio"
	"go.etcd.io/bbolt"
)

// Version of client. Set during build.
// "0.0.0" is the development version.
var Version = "0.0.0"

// TODO Reconsiliation tests
// TODO Measure download progress
// TODO Measure upload progress
// TODO Daemon mode
// TODO HTTP API

var (
	versionFlag = flag.Bool("version", false, "print program version")
	debugFlag   = flag.Bool("debug", false, "enable debug logs")
	configFlag  = flag.String("config", "", "path of config file")
	username    = flag.String("username", "", "put.io account username")
	password    = flag.String("password", "", "put.io account password")
)

var (
	configPath     string
	config         Config
	db             *bbolt.DB
	token          string
	client         *putio.Client
	localPath      string
	remoteFolderID int64
	dirCache       = NewDirCache()
)

var bucketFiles = []byte("files")

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
	log.Infof("Using config file %q", configPath)
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
	log.Infof("Using database file %q", dbPath)
	db, err = bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketFiles)
		return err
	})
	if err != nil {
		log.Fatal(err)
	}
	err = syncOnce(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	log.Infoln("Sync finished")
}

func syncOnce(ctx context.Context) error {
	err := authenticate(ctx)
	if err != nil {
		return err
	}
	err = ensureRoots(ctx)
	if err != nil {
		return err
	}
	return syncRoots(ctx)
}
