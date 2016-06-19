package glog

import (
	"os"

	"github.com/op/go-logging"
	//"github.com/gin-gonic/gin"
)

// GetLogger returns the requested logger
func GetLogger(kind string) *logging.Logger {

	var log = logging.MustGetLogger("example")
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} > %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)

	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	logging.SetBackend(backend2Formatter)

	//log.Debug(kind, " kind of logger is requested")
	return log
}
