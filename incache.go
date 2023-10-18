package incache

import (
	"sync"
	"time"
)

// Cache is a synchronised map of items that are automatically removed
// when they expire.
//
// Example:
//
// cache := incache.New(
// 	incache.WithTTL(5*time.Minute),
// 	incache.WithCleanupInterval(5*time.Minute),
// 	incache.WithMetrics(),
// )
type Cache struct {
	mu               sync.RWMutex
	items            map[string]Item
	expirationsQueue expirationsQueue
	cleaner          *cleaner

	config  Config
	metrics metrics
}

// New creates new instance of the cache.
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
	}

	if config.cleanupInterval > 0 {
		cache.cleaner = newCleaner(config.cleanupInterval)
		cache.cleaner.start(cache)
	}

	return cache
}

// Close stops the automatic cleanup process.
//
// It's not necessary to run this function if you have cleanupInterval <= 0,
// since cleaner wasn't run.
func (c *Cache) Close() {
	if c.cleaner != nil {
		c.cleaner.stop()
	}
}

// Set sets the key to hold a value.
// If key already holds a value, It will be overwritten.
func (c *Cache) Set(key string, value interface{}) {
	ttl := c.config.ttl

	c.set(key, value, ttl)
}

// SetWithTTL works similar to Set method, but with an opportunity to
// adjust a ttl for that particular key manually.
func (c *Cache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.set(key, value, ttl)
}

// SetGet sets the key to hold a value, and then returns it.
func (c *Cache) SetGet(key string, value interface{}) interface{} {
	ttl := c.config.ttl

	c.set(key, value, ttl)
	v := c.get(key)

	return v
}

// SetGetWithTTL works similar to SetGet method, but with an opportunity
// to adjust  ttl for that particular key manually.
func (c *Cache) SetGetWithTTL(key string, value interface{}, ttl time.Duration) interface{} {
	c.set(key, value, ttl)
	v := c.get(key)

	return v
}

// Get returns the value of key.
// If the key doesn't exist, nil value will be returned.
func (c *Cache) Get(key string) interface{} {
	return c.get(key)
}

// GetMultiple returns the values of all specified keys.
// For every specified key that doesn't exist, nil value will be returned.
func (c *Cache) GetMultiple(keys []string) []interface{} {
	values := make([]interface{}, len(keys))

	for i, key := range keys {
		values[i] = c.get(key)
	}

	return values
}

// GetSet returns the old value stored by key and set the new one for that key.
// If the key doesn't exist, nil value will be returned.
func (c *Cache) GetSet(key string, value interface{}) interface{} {
	ttl := c.config.ttl

	v := c.get(key)
	c.set(key, value, ttl)

	return v
}

// GetSetWithTTL works similar to GetSet method, but with an opportunity to
// adjust a ttl for that particular key manually.
func (c *Cache) GetSetWithTTL(key string, value interface{}, ttl time.Duration) interface{} {
	v := c.get(key)
	c.set(key, value, ttl)

	return v
}

// GetDelete returns the value of key and delete it.
// If the key doesn't exist, nil value will be returned.
func (c *Cache) GetDelete(key string) interface{} {
	value := c.get(key)
	if value != nil {
		c.evict(key)
	}

	return value
}

// Delete deletes the value of key.
// If the key doesn't exist, nothing will happen.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	_, ok := c.items[key]
	c.mu.Unlock()

	if ok {
		c.evict(key)
	}
}

// DeleteAll deletes all values stored in the cache.
func (c *Cache) DeleteAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]Item)
}

// DeleteExpired deletes all expired items from the cache.
// You don't need to perform it manually unless cleanupInterval is <= 0.
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

// Keys returns slice of all existing keys in the cache.
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

// Len returns the number of stored elements in the cache.
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	count := len(c.items)
	return count
}

// Has checks if the key exists in the cache.
func (c *Cache) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.items[key]
	return ok
}

// Metrics returns collected cache metrics.
func (c *Cache) Metrics() metrics {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.metrics
}

// ResetMetrics resets cache metrics.
func (c *Cache) ResetMetrics() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metrics = metrics{}
}

func (c *Cache) set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item := newItem(value, ttl)
	c.items[key] = item

	if item.CanExpire() {
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
