package limiter

type UltimateLimiter struct {
	sl *SpeedLimiter
	tl *TotalLimiter
}

func NewUltimateLimiter(speedLimit int, totalLimit int) *UltimateLimiter {
	return &UltimateLimiter{NewSpeedLimiter(speedLimit), NewTotalLimiter(totalLimit)}
}

func (l *UltimateLimiter) Recieve(addr string) bool {
	if !l.tl.Recieve(addr) {
		return false
	}

	return l.sl.Recieve(addr)
}

func (l *UltimateLimiter) End(addr string) {
	l.tl.End(addr)
}
