package storage

import "github.com/smakimka/tack/internal/model"

type Storage interface {
	Update(data model.SendData)

	CheckAlive()

	GetTacks() []string
	GetTack(string) (model.Tack, bool)
	GetProxy(tackName string, proxyName string) (model.EndpointData, bool)
	GetBalancer(tackName string, balancerName string) (model.EndpointData, bool)
}
