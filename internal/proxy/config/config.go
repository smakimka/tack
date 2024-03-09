package config

import (
	"flag"

	"github.com/BurntSushi/toml"

	"github.com/smakimka/tack/internal/model"
)

type Config struct {
	Name         string `toml:"name"`
	Logginglevel string `toml:"logging_level"`
	Monitor      bool   `toml:"monitor"`
	MonitorAddr  string `toml:"monitor_addr"`
	Endpoints    map[string]model.ProxyEndpoint
	Balancers    map[string]model.Balancer
}

func Read() (*Config, error) {
	var flagConfigFile string
	flag.StringVar(&flagConfigFile, "f", "config", "path to a config file")
	flag.Parse()

	config := &Config{}
	_, err := toml.DecodeFile(flagConfigFile, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
