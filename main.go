package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"

	"github.com/putdotio/putio-sync/http"
	"github.com/putdotio/putio-sync/sync"
)

func main() {
	log.SetFlags(0)

	// flags
	var (
		serverFlag = flag.Bool("server", false, "Run in server mode")
		debugFlag  = flag.Bool("debug", false, "Run in debug mode")
	)
	flag.Parse()

	u, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}

	err = os.MkdirAll(filepath.Join(u.HomeDir, ".putio-sync/"), 0755)
	if err != nil {
		log.Fatalln(err)
	}

	cfgpath := filepath.Join(u.HomeDir, ".putio-sync/putio-sync.db")
	sync, err := sync.NewClient(cfgpath, *debugFlag)
	if err != nil {
		log.Fatalln(err)
	}

	var server *http.Server
	if *serverFlag {
		server = http.NewServer(sync)
		err := server.Open()
		if err != nil {
			log.Fatalln(err)
		}

		go func() {
			log.Printf("Visit 'http://127.0.0.1%v'\n", server.Addr)
			log.Fatalln(server.Serve())
		}()
	} else {
		err = sync.Run()
		if err != nil {
			log.Fatalln(err)
		}
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, os.Kill)

	sig := <-sigCh
	log.Printf("%q signal received, closing...\n", sig)

	err = sync.Close()
	if err != nil {
		log.Fatalln(err)
	}

	if *serverFlag {
		err := server.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}
}
