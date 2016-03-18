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

// Config the configuration struct
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

// ReadConfig reads the json config file and adds the standard values
func ReadConfig() Config {
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

	var conf Config
	if err := json.Unmarshal(cfgData, &conf); err != nil {
		panic(err)
	}

	conf.ChangeTesholdPerc = 5 // hc value - do not expose
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
					ret = append(ret, ipnet)
				}
			}
		}
	}
	return ret, nil
}
