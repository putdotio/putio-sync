package sync

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

// Logger represents a custom log.Logger with debug methods.
type Logger struct {
	debug bool
	*log.Logger
	w io.WriteCloser
}

// NewLogger creates a new Logger. If path is not empty, it creates a log file.
func NewLogger(prefix string, debug bool, path string) *Logger {
	var w io.WriteCloser

	if path == "" {
		w = os.Stdout
	} else {
		path = filepath.Join(path, "putio-sync.log")
		f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			w = os.Stdout
		} else {
			w = f
		}
	}

	return &Logger{
		debug:  debug,
		Logger: log.New(w, prefix, log.Lshortfile|log.LstdFlags),
		w:      w,
	}
}

// Debugf calls log.Printf if debug is enabled.
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.debug {
		format += "[DEBUG] "
		l.Printf(format, v...)
	}
}

// Close closes the underlying file descriptor.
func (l *Logger) Close() error {
	return l.w.Close()
}
