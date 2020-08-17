package main

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

type Server struct {
	srv  *http.Server
	done chan struct{}
}

func NewServer(addr string) *Server {
	m := http.NewServeMux()
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte("putio-sync " + Version)) })
	s := &Server{
		srv: &http.Server{
			Addr:         addr,
			Handler:      m,
			ReadTimeout:  serverReadTimeout,
			WriteTimeout: serverWriteTimeout,
		},
		done: make(chan struct{}),
	}
	return s
}

func (s *Server) Start() {
	go func() {
		if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}
	close(s.done)
	return nil
}

func (s *Server) Wait() {
	<-s.done
}
