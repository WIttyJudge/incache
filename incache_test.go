package incache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefault(t *testing.T) {
	cache := New()
	defaultConfig := defaultConfig()

	require.NotNil(t, cache)
	assert.Equal(t, defaultConfig.ttl, cache.config.ttl)
	assert.Equal(t, defaultConfig.cleanupInterval, cache.config.cleanupInterval)
	assert.Equal(t, defaultConfig.enableMetrics, cache.config.enableMetrics)
	assert.Equal(t, defaultConfig.enableDebug, cache.config.enableDebug)
	assert.NotNil(t, cache.config.debugf)

	assert.NotNil(t, cache.items)
	assert.NotNil(t, cache.expirationsQueue)
	assert.NotNil(t, cache.metrics)
	assert.NotNil(t, cache.cleaner)
}

func TestNewCustomConfig(t *testing.T) {
	customConfig := Config{
		ttl:             1 * time.Minute,
		cleanupInterval: 1 * time.Minute,
		enableMetrics:   true,
		enableDebug:     true,
	}

	cache := New(
		WithTTL(customConfig.ttl),
		WithCleanupInterval(customConfig.cleanupInterval),
		WithMetrics(),
		WithDebug(),
	)

	require.NotNil(t, cache)
	assert.Equal(t, customConfig.ttl, cache.config.ttl)
	assert.Equal(t, customConfig.cleanupInterval, cache.config.cleanupInterval)
	assert.Equal(t, customConfig.enableMetrics, cache.config.enableMetrics)
	assert.Equal(t, customConfig.enableDebug, cache.config.enableDebug)
	assert.NotNil(t, cache.config.debugf)
}

func TestNewCustomConfigWithoutCleanup(t *testing.T) {
	cache := New(WithCleanupInterval(0))

	require.NotNil(t, cache)

	assert.Nil(t, cache.cleaner)
	assert.EqualValues(t, 0, cache.config.cleanupInterval)
}

func TestSet(t *testing.T) {
	cache := New()

	cache.Set("key1", "value1")
	assert.Equal(t, "value1", cache.Get("key1"))
	assert.Equal(t, "value1", cache.items["key1"].Value)
}

func TestSetWithTTL(t *testing.T) {
	cache := New()

	cache.SetWithTTL("key1", "value1", 5*time.Millisecond)
	assert.Equal(t, "value1", cache.Get("key1"))
	time.Sleep(5 * time.Millisecond)
	assert.Nil(t, cache.Get("key1"))
}

func TestSetGet(t *testing.T) {
	cache := New()

	recevedValue := cache.SetGet("key1", "value1")
	assert.Equal(t, "value1", recevedValue)
}

func TestSetGetWithTTL(t *testing.T) {
	cache := New()

	recevedValue := cache.SetGetWithTTL("key1", "value1", 5*time.Millisecond)
	assert.Equal(t, "value1", recevedValue)
	time.Sleep(5 * time.Millisecond)
	assert.Nil(t, cache.Get("key1"))
}

func TestGet(t *testing.T) {
	cache := New()

	cache.Set("key1", "value1")
	cache.SetWithTTL("key2", "value2", 1*time.Millisecond)

	assert.Equal(t, "value1", cache.Get("key1"))
	time.Sleep(1 * time.Millisecond)
	assert.Nil(t, cache.Get("key2"))
}

func TestGetMultiple(t *testing.T) {
	cache := New()

	cache.Set("key1", "value1")
	cache.SetWithTTL("key2", "value2", 1*time.Millisecond)

	assert.Equal(t, []interface{}{"value1"}, cache.GetMultiple([]string{"key1"}))
	time.Sleep(1 * time.Millisecond)
	assert.Equal(t, []interface{}{"value1", nil}, cache.GetMultiple([]string{"key1", "key2"}))
}

func TestGetSet(t *testing.T) {
	cache := New()

	receivedValue := cache.GetSet("key1", "value1")
	assert.Equal(t, nil, receivedValue)

	receivedValue = cache.GetSet("key1", "value2")
	assert.Equal(t, "value1", receivedValue)
}

func TestGetSetWithTTL(t *testing.T) {
	cache := New()

	receivedValue := cache.GetSetWithTTL("key1", "value1", 1*time.Minute)
	assert.Equal(t, nil, receivedValue)

	receivedValue = cache.GetSetWithTTL("key1", "value2", 5*time.Millisecond)
	assert.Equal(t, "value1", receivedValue)
	time.Sleep(5 * time.Millisecond)
	assert.Nil(t, cache.Get("key1"))
}

func TestGetDelete(t *testing.T) {
	cache := New()

	cache.Set("key1", "value1")

	recevedValue := cache.GetDelete("key1")
	assert.Equal(t, "value1", recevedValue)
	assert.Nil(t, cache.items["key1"].Value)
	assert.Nil(t, cache.GetDelete("key2"))
}

func TestDelete(t *testing.T) {
	cache := New()

	cache.Set("key1", "value1")
	assert.Equal(t, "value1", cache.items["key1"].Value)
	assert.Len(t, cache.items, 1)

	cache.GetDelete("key1")
	assert.Nil(t, cache.items["key1"].Value)
	assert.Len(t, cache.items, 0)
}

func TestDeleteAll(t *testing.T) {
	cache := New()

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	cache.DeleteAll()

	assert.Empty(t, cache.items)
}

func TestDeleteExpired(t *testing.T) {
	cache := New(WithTTL(1 * time.Millisecond))

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	time.Sleep(1 * time.Millisecond)

	cache.DeleteExpired()

	assert.Empty(t, cache.items)
}

func TestKeys(t *testing.T) {
	cache := New()

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	assert.ElementsMatch(t, []string{"key1", "key2"}, cache.Keys())
}

func TestLen(t *testing.T) {
	cache := New()

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	assert.Len(t, cache.items, 2)
	assert.Equal(t, 2, cache.Len())
}

func TestHas(t *testing.T) {
	cache := New()

	cache.Set("key1", "value1")

	assert.True(t, cache.Has("key1"))
	assert.False(t, cache.Has("nokey"))
}

func TestMetrics(t *testing.T) {
	cache := New(WithMetrics())

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Get("key1")
	cache.Get("key3")
	cache.Delete("key1")

	metrics := cache.Metrics()

	assert.EqualValues(t, 2, metrics.Insertions())
	assert.EqualValues(t, 1, metrics.Hits())
	assert.EqualValues(t, 1, metrics.Misses())
	assert.EqualValues(t, 1, metrics.Evictions())
}

func TestNoMetricsType(t *testing.T) {
	cache := New()
	metrics := cache.Metrics()

	assert.IsType(t, &noMetrics{}, metrics)
}

func TestRealMetricsType(t *testing.T) {
	cache := New(WithMetrics())
	metrics := cache.Metrics()

	assert.IsType(t, &realMetrics{}, metrics)
}

func TestResetMetrics(t *testing.T) {
	cache := New(WithMetrics())
	metrics := cache.Metrics()

	cache.Set("key1", "value1")

	assert.EqualValues(t, 1, metrics.Insertions())

	cache.ResetMetrics()

	assert.EqualValues(t, 0, metrics.Insertions())
}

func TestOnInsertion(t *testing.T) {
	cache := New()
	checkCh := make(chan struct{}, 2)

	cache.OnInsertion(func(_ string, _ interface{}) {
		checkCh <- struct{}{}
	})

	cache.Set("key1", "value1")

	assert.Eventually(t, func() bool {
		return len(checkCh) == 1
	}, time.Millisecond*500, time.Millisecond*250)
}

func TestOnEviction(t *testing.T) {
	cache := New()
	checkCh := make(chan struct{}, 1)

	cache.OnEviction(func(_ string, _ interface{}) {
		checkCh <- struct{}{}
	})

	cache.Set("key1", "value1")
	cache.Delete("key1")

	assert.Eventually(t, func() bool {
		return len(checkCh) == 1
	}, time.Millisecond*500, time.Millisecond*250)
}
