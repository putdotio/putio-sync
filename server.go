package putiosync

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/cenkalti/log"
)

const (
	serverReadTimeout     = 5 * time.Second
	serverWriteTimeout    = 10 * time.Second
	serverShutdownTimeout = 5 * time.Second
)

type httpServer struct {
	srv *http.Server
}

func newServer(addr string) *httpServer {
	m := http.NewServeMux()
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte("putio-sync")) })
	m.HandleFunc("/syncing", func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(fmt.Sprintf("%v", syncing))) })
	m.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		b, _ := json.Marshal(map[string]string{"status": syncStatus})
		_, _ = w.Write(b)
	})
	s := &httpServer{
		srv: &http.Server{
			Addr:         addr,
			Handler:      m,
			ReadTimeout:  serverReadTimeout,
			WriteTimeout: serverWriteTimeout,
		},
	}
	return s
}

func (s *httpServer) Close() {
	s.srv.Close()
}

func (s *httpServer) Start() {
	l, err := net.Listen("tcp4", s.srv.Addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Infoln("Server is listening on", l.Addr().String())
	go func() {
		if err := s.srv.Serve(l); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
}

func (s *httpServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()
	return s.srv.Shutdown(ctx)
}
