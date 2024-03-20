package balancers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundRobinSplitter(t *testing.T) {
	splitter := NewRoundRobinSplitter()
	splitter.addrs = []string{"localhost:8000", "localhost:8001", "localhost:8002"}

	assert.Equal(t, "localhost:8001", splitter.Next())
	assert.Equal(t, "localhost:8002", splitter.Next())
	assert.Equal(t, "localhost:8000", splitter.Next())
	assert.Equal(t, "localhost:8001", splitter.Next())
}
