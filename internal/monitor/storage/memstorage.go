package storage

import (
	"sync"
	"time"

	"github.com/smakimka/tack/internal/model"
)

type MemStorage struct {
	m     sync.RWMutex
	tacks map[string]model.Tack
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		m:     sync.RWMutex{},
		tacks: map[string]model.Tack{},
	}
}

func (s *MemStorage) GetTacks() []string {
	s.m.RLock()
	defer s.m.RUnlock()

	tacks := []string{}
	for tackName := range s.tacks {
		tacks = append(tacks, tackName)
	}

	return tacks
}

func (s *MemStorage) GetTack(tackName string) (model.Tack, bool) {
	s.m.RLock()
	defer s.m.RUnlock()

	tack, ok := s.tacks[tackName]

	return tack, ok
}

func (s *MemStorage) GetProxy(tackName string, proxyName string) (model.EndpointData, bool) {
	s.m.Lock()
	defer s.m.Unlock()

	var proxy model.EndpointData

	tack, ok := s.tacks[tackName]
	if !ok {
		return proxy, false
	}

	proxy, ok = tack.Proxies[proxyName]
	if !ok {
		return proxy, false
	}

	return proxy, ok
}

func (s *MemStorage) GetBalancer(tackName string, balancerName string) (model.EndpointData, bool) {
	s.m.Lock()
	defer s.m.Unlock()

	var balancer model.EndpointData

	tack, ok := s.tacks[tackName]
	if !ok {
		return balancer, false
	}

	balancer, ok = tack.Proxies[balancerName]
	if !ok {
		return balancer, false
	}

	return balancer, ok
}

func (s *MemStorage) Update(data model.SendData) {
	s.m.Lock()
	defer s.m.Unlock()

	s.tacks[data.Name] = model.Tack{
		Alive:       true,
		Proxies:     data.Endpoints,
		Balancers:   data.Balancers,
		LastUpdated: time.Now(),
	}
}

func (s *MemStorage) CheckAlive() {
	s.m.Lock()
	defer s.m.Unlock()

	now := time.Now()
	for tackName, tack := range s.tacks {
		if now.Sub(tack.LastUpdated).Seconds() > 10 {
			tack.Alive = false
			s.tacks[tackName] = tack
		}
	}
}
