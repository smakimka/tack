package model

type Endpoint struct {
	Addr       string `toml:"addr"`
	Workers    int    `toml:"workers"`
	SpeedLimit int    `toml:"speed_limit"`
	TotalLimit int    `toml:"total_limit"`
}

type ProxyEndpoint struct {
	Endpoint
}

type Balancer struct {
	Endpoint
	Type  string   `toml:"type"`
	Addrs []string `toml:"addrs"`
}
