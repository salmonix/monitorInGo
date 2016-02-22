package logger

import (
	"bytes"
	"log"
)

// GetLogger returns the standard logger
func GetLogger() *log.Logger {
	var buf bytes.Buffer
	logger := log.New(&buf, "logger: ", log.Lshortfile)
	return logger
}
