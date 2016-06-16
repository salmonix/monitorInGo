package main

import (
	"fmt"
	"gmon/watch"
	c "gmon/watch/config"
	"gmon/watch/rest"
	"os"
	"strconv"
	"time"

	"github.com/op/go-logging"
)

var l = getLog()

func main() {
	conf := c.GetConfig()
	fmt.Println(conf)
	l.Info("Starting to poll")
	watcher := startPolling(conf)
	l.Info("Getting to router")
	router := rest.GetRouter(watcher, &conf)
	router.Run(":" + strconv.Itoa(conf.Port))
}

// StartPolling starts the polling loop
func startPolling(conf c.Config) *watch.WatchingContainer {

	// add a lock and test it
	// add SIG handling: SIGTERM, SIGHUP
	l.Info("Asking for new container")
	p := watch.NewContainer(conf.ChangeTesholdPerc)
	l.Info("New container is made")
	go func() {
		for {
			p.Refresh()
			time.Sleep(time.Duration(conf.ScanIntervalSec))
		}
	}()
	return p
}

func getLog() *logging.Logger {

	var log = logging.MustGetLogger("example")
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} > %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)

	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	logging.SetBackend(backend2Formatter)

	return log
}
