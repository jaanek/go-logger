package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

// Logger
type Logger struct {
	*log.Logger
	shutdown chan struct{}
}

func (l *Logger) Close() {
	l.shutdown <- struct{}{}
}

var DefaultLogger, _ = New("")

func New(filepath string) (*Logger, error) {
	if filepath == "" {
		return &Logger{
			Logger:   log.New(os.Stdout, "", log.LstdFlags),
			shutdown: make(chan struct{}),
		}, nil
	}
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return &Logger{}, fmt.Errorf("Cannot open log file: %v, %v", filepath, err)
	}
	logger := &Logger{
		Logger:   log.New(io.MultiWriter(file, os.Stdout), "", log.LstdFlags),
		shutdown: make(chan struct{}),
	}
	go func() {
		select {
		case <-logger.shutdown:
			file.Close()
		}
	}()
	return logger, nil
}
