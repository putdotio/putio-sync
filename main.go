package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/adrg/xdg"
	"github.com/cenkalti/log"
	"github.com/putdotio/go-putio"
	"go.etcd.io/bbolt"
)

// Version of client. Set during build.
// "0.0.0" is the development version.
var Version = "0.0.0"

// TODO HTTP API
// TODO websocket endpoint for progress updates
// TODO listen websocket for remote events
// TODO watch local fs for changes

var (
	versionFlag = flag.Bool("version", false, "print program version")
	debugFlag   = flag.Bool("debug", false, "enable debug logs")
	configFlag  = flag.String("config", "", "path of config file")
	username    = flag.String("username", "", "put.io account username")
	password    = flag.String("password", "", "put.io account password")
	dryrun      = flag.Bool("dryrun", false, "do not make changes on filesystems")
	repeat      = flag.Duration("repeat", 0, "sync repeatedly, pause given duration between syncs")
	server      = flag.String("server", "", "listen address for HTTP API")
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
	var srv *Server
	if *server != "" {
		srv = NewServer(*server)
		srv.Start()
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-ch
		log.Noticef("Received %s. Stopping sync.", sig)
		cancel()

		if srv != nil {
			if err := srv.Shutdown(); err != nil {
				log.Fatalln("Server shutdown failed:", err)
			}
		}
	}()
REPEAT_LOOP:
	for {
		err = syncOnce(ctx)
		if err != nil {
			if *repeat == 0 {
				log.Fatal(err)
			}
			log.Error(err)
		} else {
			log.Infoln("Sync finished successfully")
		}
		if *repeat == 0 {
			break
		}
		select {
		case <-time.After(*repeat):
		case <-ctx.Done():
			break REPEAT_LOOP
		}
	}
	if srv != nil {
		srv.Wait()
		log.Debug("Server has shutdown successfully")
	}
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
