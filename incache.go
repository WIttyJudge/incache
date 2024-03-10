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
//
//	incache.WithTTL(5*time.Minute),
//	incache.WithCleanupInterval(5*time.Minute),
//	incache.WithMetrics(),
//
// )
type Cache struct {
	mu               sync.RWMutex
	items            map[string]Item
	expirationsQueue expirationsQueue
	cleaner          *cleaner
	eventHandlers    *eventHandlers

	config  Config
	metrics metrics
}

// New creates new instance of the cache.
func New(conf ...configFunc) *Cache {
	config := defaultConfig()
	for _, fn := range conf {
		fn(&config)
	}

	if !config.enableDebug {
		config.debugf = func(format string, v ...any) {}
	}

	cache := &Cache{
		mu:               sync.RWMutex{},
		items:            make(map[string]Item),
		expirationsQueue: make(map[string]time.Time),
		eventHandlers:    newEventHandlers(),

		config:  config,
		metrics: newNoMetrics(),
	}

	if config.enableMetrics {
		cache.metrics = newRealMetrics()
	}

	if config.cleanupInterval > 0 {
		cache.cleaner = newCleaner(config.cleanupInterval)
		cache.cleaner.start(cache)
	}

	return cache
}

// Close allows you to stop automatic cleaner manually and wait for the the
// exeuction of all events.
//
// There is no needs to run this function, if you don't use event and there is a
// cleanupInterval <= 0, since cleaner whouldn't be run in this case.
func (c *Cache) Close() {
	if c.cleaner != nil {
		c.config.debugf("[close] closing cleaner")
		c.cleaner.close()
	}

	c.config.debugf("[close] waiting for the execution of all events")
	c.eventHandlers.Wait()
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

	c.metrics.reset()
}

func (c *Cache) OnInsertion(fn func(key string, value interface{})) {
	c.eventHandlers.OnInsertion(fn)
}

func (c *Cache) OnEviction(fn func(key string, value interface{})) {
	c.eventHandlers.OnEviction(fn)
}

func (c *Cache) set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.eventHandlers.onInsertion(key, value)

	item := newItem(value, ttl)
	c.items[key] = item

	if item.CanExpire() {
		c.expirationsQueue[key] = item.ExpiresAt
	}

	c.config.debugf("[set] key: '%s', item: %+v", key, item)

	c.metrics.incrementInsertions()
}

func (c *Cache) get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item := c.items[key]
	value := item.Value

	if value == nil {
		c.config.debugf("[get] no value was found for the key: '%s'", key)

		c.metrics.incrementMisses()
		return nil
	}

	if item.Expired() {
		c.config.debugf("[get] received value for the key: '%s' is expired", key)

		c.metrics.incrementMisses()
		return nil
	}

	c.metrics.incrementHits()

	c.config.debugf("[get] key: '%s', value: %+v", key, value)

	return value
}

func (c *Cache) evict(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	value := c.items[key]
	c.eventHandlers.onEviction(key, value)

	delete(c.items, key)
	delete(c.expirationsQueue, key)

	c.config.debugf("[evict] key: '%s'", key)
	c.metrics.incrementEvictions()
}
