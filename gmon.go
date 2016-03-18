package main

import (
	"gmon/watch"
	c "gmon/watch/config"
	"gmon/watch/rest"
	"log"
	"strconv"
	"time"
)

var l *log.Logger

func main() {
	conf := c.ReadConfig()
	watcher := startPolling(conf)
	router := rest.GetRouter(watcher, &conf)
	router.Run(":" + strconv.Itoa(conf.Port))
}

// StartPolling starts the polling loop
func startPolling(conf c.Config) *watch.WatchingContainer {

	// add a lock and test for it
	p := watch.NewContainer(conf.ChangeTesholdPerc)
	go func() {
		for {
			p.Refresh()
			time.Sleep(time.Duration(conf.ScanIntervalSec))
		}
	}()
	return p
}
