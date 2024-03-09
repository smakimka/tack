package limiter

import (
	"sync"
)

type TotalLimiter struct {
	m           sync.Mutex
	limit       int
	clientCount map[string]int
}

func NewTotalLimiter(limit int) *TotalLimiter {
	return &TotalLimiter{m: sync.Mutex{}, limit: limit, clientCount: map[string]int{}}
}

func (l *TotalLimiter) Recieve(addr string) bool {
	l.m.Lock()
	defer l.m.Unlock()

	if l.clientCount[addr]+1 > l.limit {
		return false
	}

	l.clientCount[addr]++
	return true
}

func (l *TotalLimiter) End(addr string) {
	l.m.Lock()
	defer l.m.Unlock()

	l.clientCount[addr]--
}
