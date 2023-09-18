package cache

import (
	"sync"
)

type Cache struct {
	mu    sync.RWMutex
	items map[string]Item

	options Options
	metrics *Metrics
}

func New(opts ...optionsFunc) *Cache {
	options := DefaultOptions()
	for _, fn := range opts {
		fn(&options)
	}

	return &Cache{
		mu:    sync.RWMutex{},
		items: make(map[string]Item),

		options: options,
		metrics: NewMetrics(),
	}

}

// Set the key to hold a value.
// If key already holds a value, It will be overwritten.
func (c *Cache) Set(key string, value interface{}) {
	c.set(key, value)
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
	v := c.get(key)
	c.set(key, value)

	return v
}

// Get the value of key and delete it.
// If the key doesn't exist, nil value will be returned.
func (c *Cache) GetDelete(key string) interface{} {
	value := c.get(key)
	if value != nil {
		c.delete(key)
	}

	return value
}

// Delete the value of key.
// If the key doesn't exist, nothing will happen.
func (c *Cache) Delete(key string) {
	c.delete(key)
}

// Delete all values stored in the cache.
func (c *Cache) DeleteAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]Item)
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

// Get count of keys in the cache.
func (c *Cache) KeysCount() int {
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

// Get pointer to Metrics structure that collects an important metrics during
// the work with cache.
func (c *Cache) Metrics() *Metrics {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.metrics
}

func (c *Cache) set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item := Item{
		Value: value,
	}
	c.items[key] = item

	if c.options.enableMetrics {
		c.metrics.incrInsertions()
	}
}

func (c *Cache) get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value := c.items[key].Value
	if value == nil {
		if c.options.enableMetrics {
			c.metrics.incrMisses()
		}
		return nil
	}

	if c.options.enableMetrics {
		c.metrics.incrHits()
	}

	return value
}

func (c *Cache) delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}
