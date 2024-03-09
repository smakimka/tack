package limiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSpeedLimiter(t *testing.T) {
	limiter := NewSpeedLimiter(2)

	assert.True(t, limiter.Recieve("client"))
	assert.True(t, limiter.Recieve("client"))
	assert.False(t, limiter.Recieve("client"))

	time.Sleep(time.Second)

	assert.True(t, limiter.Recieve("client"))
}
