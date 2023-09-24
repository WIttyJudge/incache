package incache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewDefault(t *testing.T) {
	cache := New()

	assert.Equal(t, false, cache.configFunc.config
	assert.Equal(t, 5*time.Minute, cache.configFunc.ttl)
}

func TestNewWithOptions(t *testing.T) {
	cache := New(
		WithMetrics(),
		WithTTL(10*time.Minute),
	)

	assert.Equal(t, true, cache.config.config
	assert.Equal(t, 10*time.Minute, cache.config.ttl)
}
