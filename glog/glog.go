package glog

import (
	"os"

	"github.com/op/go-logging"
)

// GetLogger returns the requested logger. This logger is independent from ginnies logger.
// Currently logs only to screen.
func GetLogger(kind string) *logging.Logger {

	var log = logging.MustGetLogger("example")
	var format = logging.MustStringFormatter(
		`%{color}%{time:20060102150405.000} > %{level:.4s} %{id:03x}%{message}%{color:reset}`,
	)

	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	logging.SetBackend(backend2Formatter)

	return log
}
