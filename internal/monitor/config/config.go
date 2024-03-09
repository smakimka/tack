package config

import "flag"

type Config struct {
	Addr         string
	Logginglevel string
}

func Read() *Config {
	var addr string
	var loggingLevel string
	flag.StringVar(&addr, "a", "localhost:8092", "addr of service")
	flag.StringVar(&loggingLevel, "l", "error", "logging level")

	flag.Parse()

	return &Config{
		Addr:         addr,
		Logginglevel: loggingLevel,
	}
}
