package cache

import "time"

type Config struct {
	ttl           time.Duration
	enableMetrics bool
}

type configFunc func(*Config)

func defaultConfig() Config {
	return Config{
		ttl:           5 * time.Minute,
		enableMetrics: false,
	}
}

// Sets default TTL for all items that whould be stored in the cache.
func WithTTL(ttl time.Duration) configFunc {
	return func(conf *Config) {
		conf.ttl = ttl
	}
}

// Enables the collection of metrics that run throughout the cache work.
func WithMetrics() configFunc {
	return func(conf *Config) {
		conf.enableMetrics = true
	}
}
