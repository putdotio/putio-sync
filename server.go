package putiosync

import (
	"context"
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
	go func() {
		if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
}

func (s *httpServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()
	return s.srv.Shutdown(ctx)
}
