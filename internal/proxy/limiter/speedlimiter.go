package limiter

import (
	"time"
)

type SpeedLimiter struct {
	limit      int
	clientData map[string]clientLimitData
}

type clientLimitData struct {
	count int
	start time.Time
}

func NewSpeedLimiter(limit int) *SpeedLimiter {
	return &SpeedLimiter{limit: limit, clientData: map[string]clientLimitData{}}
}

func (l *SpeedLimiter) Recieve(addr string) bool {
	clientData, ok := l.clientData[addr]
	if !ok {
		clientData = clientLimitData{1, time.Now()}
		l.clientData[addr] = clientData
		return true
	}

	clientData.count++
	now := time.Now()
	if clientData.count > l.limit && now.Sub(clientData.start) <= time.Second {
		return false
	}

	if now.Sub(clientData.start) > time.Second {
		clientData.start = now
		clientData.count = 1
	}

	l.clientData[addr] = clientData
	return true
}

func (l *SpeedLimiter) End(addr string) {}
