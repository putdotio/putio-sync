package sync

import (
	"log"
	"os"
)

type Logger struct {
	debug bool
	*log.Logger
}

func NewLogger(prefix string, debug bool) *Logger {
	return &Logger{
		debug:  debug,
		Logger: log.New(os.Stdout, prefix, 0),
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.debug {
		l.Printf(format, v...)
	}
}

func (l *Logger) Debugln(v ...interface{}) {
	if l.debug {
		l.Println(v...)
	}
}
