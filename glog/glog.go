package glog

import (
	"os"

	"github.com/op/go-logging"
	"golang.org/x/crypto/ssh/terminal"
)

// NOTE: This pacakge is only to make a preliminary interface for global logging later.

// L represents the global alias for the logger
var L = logging.MustGetLogger("gmon")

// GetLogger returns the configured logger.
// This function is a thunk as it should implement various back-ends depending on the
// environment.
func init() {

	// for interactive terminal we write coloring
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		format = logging.MustStringFormatter(
			`be1 :%{color}%{time:15:04:05.000} %{shortfunc} %{level:.4s} %{id:03x}%{color:reset} %{message}`,
		)
	}

	var format = logging.MustStringFormatter(
		`be1 :%{color}%{time:15:04:05.000} %{shortfunc} %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)

	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")
	logging.SetBackend(backend1Leveled, backend2Formatter)

}

// GetLogger returns the global logger instance
func GetLogger() *logging.Logger {
	return L
}
