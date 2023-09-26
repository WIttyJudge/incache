package incache

import "time"

// Config for cache.
type Config struct {
	ttl             time.Duration
	cleanupInterval time.Duration
	enableMetrics   bool
}

type configFunc func(*Config)

// DefaultConfig initializes config with default values.
func defaultConfig() Config {
	return Config{
		ttl:             5 * time.Minute,
		cleanupInterval: 5 * time.Minute,
		enableMetrics:   false,
	}
}

// Sets default TTL for all items that would be stored in the cache.
// TTL <= 0 means that the item won't have expiration time at all.
func WithTTL(ttl time.Duration) configFunc {
	return func(conf *Config) {
		conf.ttl = ttl
	}
}

// Interval between removing expired items.
// If the interval is less than or equal to 0, no automatic clearing
// is performed.
func WithCleanupInterval(interval time.Duration) configFunc {
	return func(conf *Config) {
		conf.cleanupInterval = interval
	}
}

// Enables the collection of metrics that run throughout the cache work.
func WithMetrics() configFunc {
	return func(conf *Config) {
		conf.enableMetrics = true
	}
}
