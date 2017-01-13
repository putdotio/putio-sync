package http

import (
	"net"
	"net/http"

	"github.com/putdotio/putio-sync/sync"
)

const (
	defaultAddr = ":3000"
)

// Server represent the HTTP interface to the sync client.
type Server struct {
	ln      net.Listener
	Handler *Handler
	Addr    string
}

// NewServer returns a new instance of Server.
func NewServer(sync *sync.Client) *Server {
	return &Server{
		Handler: NewHandler(sync),
		Addr:    defaultAddr,
	}
}

// Open opens up the underlying socket for the HTTP server.
func (s *Server) Open() error {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.ln = ln
	return nil
}

// Close closes the underlying socket.
func (s *Server) Close() error {
	if s.ln == nil {
		return nil
	}

	return s.ln.Close()
}

// Port returns the port that the server is open on. Only valid after open.
func (s *Server) Port() int {
	return s.ln.Addr().(*net.TCPAddr).Port
}

// Serve serves the HTTP requests.
func (s *Server) Serve() error {
	// remember the last paused state and continue syncing
	if !s.Handler.sync.Config.IsPaused {
		err := s.Handler.sync.Run()
		if err != nil {
			return err
		}
	}

	return http.Serve(s.ln, s.Handler)
}
