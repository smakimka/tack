package balancers

import (
	"net"
	"sync"
	"sync/atomic"
)

type Splitter interface {
	AddAddrs(addrs []string) error
	Next() string
}

type RoundRobinSplitter struct {
	m     sync.RWMutex
	addrs []string
	idx   atomic.Uint64
}

func (s *RoundRobinSplitter) Next() string {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.addrs[s.idx.Add(1)%uint64(len(s.addrs))]
}

func (s *RoundRobinSplitter) AddAddrs(addrs []string) error {
	for i := range addrs {
		addr, err := net.ResolveTCPAddr("tcp", addrs[i])
		if err != nil {
			return err
		}
		addrs[i] = addr.String()
	}

	return nil
}

func NewRoundRobinSplitter() *RoundRobinSplitter {
	return &RoundRobinSplitter{sync.RWMutex{}, []string{}, atomic.Uint64{}}
}
