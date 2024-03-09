package limiter

type RateLimiter interface {
	Recieve(addr string) bool
	End(addr string)
}

func New(speedLimit int, totalLimit int) RateLimiter {
	var rateLimiter RateLimiter
	if speedLimit != 0 || totalLimit != 0 {
		if speedLimit != 0 && totalLimit != 0 {
			rateLimiter = NewUltimateLimiter(speedLimit, totalLimit)
		} else {
			if totalLimit != 0 {
				rateLimiter = NewTotalLimiter(totalLimit)
			} else {
				rateLimiter = NewSpeedLimiter(speedLimit)
			}
		}
	}

	return rateLimiter
}
