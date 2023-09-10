package main

import (
	"fmt"
	"sync"
)

type Item struct {
	Value interface{}
}

type Cache struct {
	items map[string]Item
	mu    sync.RWMutex
}

func New() *Cache {
	return &Cache{
		mu:    sync.RWMutex{},
		items: make(map[string]Item),
	}
}

// Set key to hold the string value.
// If key already holds a value, I will be overwritten.
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = Item{
		Value: value,
	}
}

// Get the value of key.
// If the key doesn't exist, nil value will be returned.
func (c *Cache) Get(key string) interface{} {
	return c.get(key)
}

// Returns the values of all specified keys.
// For every specified key that doesn't exist, nil value will be returned.
func (c *Cache) MGet(keys ...string) []interface{} {
	values := make([]interface{}, len(keys))

	for i, key := range keys {
		values[i] = c.get(key)
	}

	return values
}

// Get the value of key and delete the key.
// If the key doesn't exist, nil value will be returned.
func (c *Cache) GetDel(key string) interface{} {
	value := c.get(key)
	if value != nil {
		c.Del(key)
	}

	return value
}

// Removes the specified keys.
// A key is ignored if it doesn't exist.
func (c *Cache) Del(key ...string) {
	for _, k := range key {
		delete(c.items, k)
	}
}

// Returns all stored keys.
func (c *Cache) Keys() []string {
	keys := make([]string, len(c.items))

	i := 0
	for key := range c.items {
		keys[i] = key
		i++
	}

	return keys
}

func (c *Cache) Flush() {
	c.items = make(map[string]Item)
}

func (c *Cache) get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.items[key].Value
}

func main() {
	cache := New()

	cache.Set("Test", 1)
	cache.Set("1", "STRING")
	// cache.Clear()

	result := cache.MGet("Test", "1")
	fmt.Println(result)
}
