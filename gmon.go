package main

import (
	"fmt"
	"gmon/glog"
	"gmon/watch"
	c "gmon/watch/config"
	"gmon/watch/rest"
	"strconv"
	"time"
)

// Version variables filled by the Makefile
var (
	Release  string
	Revision string
	Built    string
	Branch   string
)

// ReleaseBuild can be overwritten by ldflags. See: Makefile
var ReleaseBuild string = "false" // release or testing

var l = glog.GetLogger("main")

func main() {
	conf := c.GetConfig()
	fmt.Println(conf)
	watcher := startPolling(conf)
	router := rest.GetRouter(watcher, &conf)
	router.Run(":" + strconv.Itoa(conf.Port))
}

// StartPolling starts the polling loop
func startPolling(conf c.Config) *watch.WatchingContainer {

	// TODO: add a lock and test it
	// TODO: add SIG handling: SIGTERM, SIGHUP
	sleepTime := time.Duration(int64(conf.ScanIntervalSec)) * time.Second
	w := watch.NewContainer(conf.ChangeTesholdPerc)
	go func() {
		for {
			w.Refresh()
			time.Sleep(sleepTime)
		}
	}()
	return w
}
