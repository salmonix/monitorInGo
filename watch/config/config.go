package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"time"
)

// Config the configuration struct that is initialized at first call of GetConfig
type Config struct {
	ScanIntervalSec   int
	ChangeTesholdPerc float64
	Port              int
	Pid               int
	Hostname          string
	HostIP            []*net.IPNet
	StartTime         int32
	// Overseer []string
}

var conf Config

// GetConfig initializes the conf variable reading the JSON config file or
// returning an already existing conf global.
func GetConfig() Config {

	if conf.Pid > 1 {
		return conf
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Cannot get `pwd`: %s", err))
	}
	cwd = path.Join(cwd, "gmonrc.js")
	var cfg = flag.String("config", cwd, "config PATH_TO_CONFIG. defaults to `pwd`")
	flag.Parse()
	cfgData, err := ioutil.ReadFile(*cfg)
	if err != nil {
		panic(fmt.Errorf("Cannot read config file: %s", err))
	}

	if err := json.Unmarshal(cfgData, &conf); err != nil {
		panic(err)
	}

	if conf.ChangeTesholdPerc < 1 {
		conf.ChangeTesholdPerc = 5
	}
	conf.Pid = os.Getpid()
	conf.Hostname, _ = os.Hostname()
	conf.HostIP, _ = getIps()
	conf.StartTime = int32(time.Now().Unix())

	if conf.ScanIntervalSec == 0 {
		conf.ScanIntervalSec = 15
	}

	return conf
}

func getIps() ([]*net.IPNet, error) {

	var ret []*net.IPNet

	ifaces, err := net.Interfaces()
	if err != nil {
		return ret, err
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return ret, nil
		}
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ret = append(ret, ipnet) // FIXME: when marshalled it is doing somethign very funny for bitmask
				}
			}
		}
	}
	return ret, nil
}
