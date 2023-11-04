package incache

import (
	"log"
	"os"
	"time"
)

// Config for cache.
type Config struct {
	ttl             time.Duration
	cleanupInterval time.Duration
	enableMetrics   bool
	enableDebug     bool
	// It only works when debug if enabled
	debugf func(format string, v ...any)
}

type configFunc func(*Config)

// DefaultConfig initializes config with default values.
func defaultConfig() Config {
	return Config{
		ttl:             5 * time.Minute,
		cleanupInterval: 5 * time.Minute,
		enableMetrics:   false,
		enableDebug:     false,
		debugf:          log.New(os.Stdout, "[incache]", 0).Printf,
	}
}

// WithTTL sets the default TTL for all items that would be stored in
// the cache. TTL <= 0 means that the item won't have expiration time at all.
func WithTTL(ttl time.Duration) configFunc {
	return func(conf *Config) {
		conf.ttl = ttl
	}
}

// WithCleanupInterval sets the interval between removing expired items.
// If the interval is less than or equal to 0, no automatic clearing
// is performed.
func WithCleanupInterval(interval time.Duration) configFunc {
	return func(conf *Config) {
		conf.cleanupInterval = interval
	}
}

// WithMetrics enables the collection of metrics that run throughout
// the cache work.
func WithMetrics() configFunc {
	return func(conf *Config) {
		conf.enableMetrics = true
	}
}

// WithMetrics enables debug mode.
// Debug mode allows the caching system to log debug information.
func WithDebug() configFunc {
	return func(config *Config) {
		config.enableDebug = true
	}
}

// WithDebugf sets a custom debug log function in the configuration.
// This function is responsible for logging debug messages.
func WithDebugf(fn func(format string, v ...any)) configFunc {
	return func(config *Config) {
		config.debugf = fn
	}
}
