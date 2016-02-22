package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gmon/watch"
	"gmon/watch/rest"
	"io/ioutil"
	"os"
	"path"
	"time"
)

func main() {
	c := readConfig()
	polling := startPolling(c)
	router := rest.GetRouter(polling)
	router.Run(":8080")
}

type config struct {
	scanIntervalSec   int
	changeTesholdPerc float64
}

// read the json config file and retrun the map
func readConfig() config {
	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Cannot get `pwd`: %s", err))
	}
	cwd = path.Join(cwd, "gmon.conf")
	var cfg = flag.String("config", cwd, "config PATH_TO_CONFIG. defaults to `pwd`")
	flag.Parse()
	cfgData, err := ioutil.ReadFile(*cfg)
	if err != nil {
		panic(fmt.Errorf("Cannot read config file: %s", err))
	}

	var conf config
	if err := json.Unmarshal(cfgData, &conf); err != nil {
		panic(err)
	}
	return conf

}

// StartPolling starts the polling loop. NOTE: we need to put a lock on the container
func startPolling(c config) *watch.WatchingContainer {

	// add a lock and test for it
	p := watch.NewContainer(c.changeTesholdPerc)
	go func() {
		for {
			p.Refresh()
			time.Sleep(time.Duration(c.scanIntervalSec))
		}
	}()
	return p
}
