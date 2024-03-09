package limiter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTotalLimiter_Recieve(t *testing.T) {
	limiter := NewTotalLimiter(2)

	res := limiter.Recieve("client")
	assert.True(t, res)

	res = limiter.Recieve("client")
	assert.True(t, res)

	res = limiter.Recieve("client")
	assert.False(t, res)
}

func TestTotalLimiter_End(t *testing.T) {
	limiter := NewTotalLimiter(2)

	limiter.Recieve("client")
	limiter.Recieve("client")

	limiter.End("client")

	res := limiter.Recieve("client")
	assert.True(t, res)

	res = limiter.Recieve("client")
	assert.False(t, res)
}
