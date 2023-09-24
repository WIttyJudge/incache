package incache

import (
	"sync"
	"time"
)

// Cache is a synchronised map of items that are automatically removed
// when they expire.
type Cache struct {
	mu               sync.RWMutex
	items            map[string]Item
	expirationsQueue map[string]time.Time

	config  Config
	metrics metrics

	closeCh chan struct{}
}

// Creates new instance of the cache.
func New(conf ...configFunc) *Cache {
	config := defaultConfig()
	for _, fn := range conf {
		fn(&config)
	}

	cache := &Cache{
		mu:               sync.RWMutex{},
		items:            make(map[string]Item),
		expirationsQueue: make(map[string]time.Time),

		config:  config,
		metrics: newMetrics(),

		closeCh: make(chan struct{}),
	}

	if config.cleanupInterval > 0 {
		go cache.runAutomaticCleanup()
	}

	return cache
}

// Stops the automatic cleanup process.
// You don't need to run this function if you have cleanupInterval <= 0.
func (c *Cache) Close() {
	close(c.closeCh)
}

// Set the key to hold a value.
// If key already holds a value, It will be overwritten.
func (c *Cache) Set(key string, value interface{}) {
	ttl := c.config.ttl

	c.set(key, value, ttl)
}

// Similar to Set method, but with an opportunity to adjust a ttl for that
// particular key manually.
func (c *Cache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.set(key, value, ttl)
}

// Set the key to hold a value, and then returns it.
func (c *Cache) SetGet(key string, value interface{}) interface{} {
	ttl := c.config.ttl

	c.set(key, value, ttl)
	v := c.get(key)

	return v
}

// Similar to SetGet method, but with an opportunity to adjust  ttl for that
// particular key manually.
func (c *Cache) SetGetWithTTL(key string, value interface{}, ttl time.Duration) interface{} {
	c.set(key, value, ttl)
	v := c.get(key)

	return v
}

// Get the value of key.
// If the key doesn't exist, nil value will be returned.
func (c *Cache) Get(key string) interface{} {
	return c.get(key)
}

// Get the values of all specified keys.
// For every specified key that doesn't exist, nil value will be returned.
func (c *Cache) GetMultiple(keys []string) []interface{} {
	values := make([]interface{}, len(keys))

	for i, key := range keys {
		values[i] = c.get(key)
	}

	return values
}

// Get the old value stored by key and set the new one for that key.
// If the key doesn't exist, nil value will be returned.
func (c *Cache) GetSet(key string, value interface{}) interface{} {
	ttl := c.config.ttl

	v := c.get(key)
	c.set(key, value, ttl)

	return v
}

// Similar to GetSet method, but with an opportunity to adjust a ttl for that
// particular key manually.
func (c *Cache) GetSetWithTTL(key string, value interface{}, ttl time.Duration) interface{} {
	v := c.get(key)
	c.set(key, value, ttl)

	return v
}

// Get the value of key and delete it.
// If the key doesn't exist, nil value will be returned.
func (c *Cache) GetDelete(key string) interface{} {
	value := c.get(key)
	if value != nil {
		c.evict(key)
	}

	return value
}

// Delete the value of key.
// If the key doesn't exist, nothing will happen.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	value := c.items[key].Value
	c.mu.Unlock()

	if value != nil {
		c.evict(key)
	}
}

// Delete all values stored in the cache.
func (c *Cache) DeleteAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]Item)
}

// DeleteExpired deletes all expired items from the cache.
func (c *Cache) DeleteExpired() {
	c.mu.Lock()

	timeNow := time.Now()
	expiredKeys := make([]string, 0, len(c.expirationsQueue))

	for key, time := range c.expirationsQueue {
		if timeNow.Before(time) {
			continue
		}

		expiredKeys = append(expiredKeys, key)
	}

	c.mu.Unlock()

	for _, key := range expiredKeys {
		c.evict(key)
	}
}

// Get slice of all existing keys in the cache.
func (c *Cache) Keys() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	keys := make([]string, len(c.items))

	i := 0
	for key := range c.items {
		keys[i] = key
		i++
	}

	return keys
}

// Returns the number of stored elements in the cache.
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	count := len(c.items)
	return count
}

// Check if the key exists in the cache.
func (c *Cache) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.items[key]
	return ok
}

// Get collected cache metrics.
func (c *Cache) Metrics() metrics {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.metrics
}

// Resets cache metrics.
func (c *Cache) ResetMetrics() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metrics = metrics{}
}

// If cleanupInterval more than 0, it will run inside goroutine.
// Checks if there are expired items and deletes them.
func (c *Cache) runAutomaticCleanup() {
	ticker := time.NewTicker(c.config.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-c.closeCh:
			return
		}
	}
}

func (c *Cache) set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item := newItem(value, ttl)
	c.items[key] = item

	if !item.ExpiresAt.IsZero() {
		c.expirationsQueue[key] = item.ExpiresAt
	}

	c.metricsIncrInsertions()
}

func (c *Cache) get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item := c.items[key]
	value := item.Value

	if value == nil {
		c.metricsIncrMisses()
		return nil
	}

	if item.Expired() {
		c.metricsIncrMisses()
		return nil
	}

	c.metricsIncrHits()

	return value
}

func (c *Cache) evict(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
	delete(c.expirationsQueue, key)
	c.metricsIncrEvictions()
}

func (c *Cache) metricsIncrInsertions() {
	if c.config.enableMetrics {
		c.metrics.incrInsertions()
	}
}

func (c *Cache) metricsIncrHits() {
	if c.config.enableMetrics {
		c.metrics.incrHits()
	}
}

func (c *Cache) metricsIncrMisses() {
	if c.config.enableMetrics {
		c.metrics.incrMisses()
	}
}

func (c *Cache) metricsIncrEvictions() {
	if c.config.enableMetrics {
		c.metrics.incrEvictions()
	}
}
