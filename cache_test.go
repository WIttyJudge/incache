package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_New_Default(t *testing.T) {
	cache := New()

	assert.Equal(t, false, cache.options.enableMetrics)
	assert.Equal(t, 5*time.Minute, cache.options.ttl)
}

func Test_New_With_Options(t *testing.T) {
	cache := New(
		WithMetrics(),
		WithTTL(10*time.Minute),
	)

	assert.Equal(t, true, cache.options.enableMetrics)
	assert.Equal(t, 10*time.Minute, cache.options.ttl)
}
