package model

import (
	"time"
)

type Tack struct {
	Proxies     map[string]EndpointData
	Balancers   map[string]EndpointData
	Alive       bool
	LastUpdated time.Time
}
