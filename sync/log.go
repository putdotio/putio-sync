package sync

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

type Logger struct {
	debug bool
	*log.Logger
	w io.WriteCloser
}

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
		Logger: log.New(w, prefix, 0),
		w:      w,
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.debug {
		format += "[DEBUG] "
		l.Printf(format, v...)
	}
}

func (l *Logger) Close() error {
	return l.w.Close()
}
