package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// Config the configuration struct
type Config struct {
	ScanIntervalSec   int
	ChangeTesholdPerc float64
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

	if conf.ScanIntervalSec == 0 {
		conf.ScanIntervalSec = 15
	}

	fmt.Println(json.Marshal(conf))
	return conf
}
