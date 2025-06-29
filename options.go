package incache

import (
	"log"
	"os"
	"time"
)

// Options for cache.
type Options struct {
	ttl             time.Duration
	cleanupInterval time.Duration
	enableMetrics   bool
	enableDebug     bool
	// It only works when debug if enabled
	debugf func(format string, v ...any)
}

type optionsFunc func(*Options)

// DefaultOptions initializes options with default values.
func defaultOptions() Options {
	return Options{
		ttl:             5 * time.Minute,
		cleanupInterval: 5 * time.Minute,
		enableMetrics:   false,
		enableDebug:     false,
		debugf:          log.New(os.Stdout, "[incache]", 0).Printf,
	}
}

// WithTTL sets the default TTL for all items that would be stored in
// the cache. TTL <= 0 means that the item won't have expiration time at all.
func WithTTL(ttl time.Duration) optionsFunc {
	return func(o *Options) {
		o.ttl = ttl
	}
}

// WithCleanupInterval sets the interval between removing expired items.
// If the interval is less than or equal to 0, no automatic clearing
// is performed.
func WithCleanupInterval(interval time.Duration) optionsFunc {
	return func(o *Options) {
		o.cleanupInterval = interval
	}
}

// WithMetrics enables the collection of metrics that run throughout
// the cache work.
func WithMetrics() optionsFunc {
	return func(o *Options) {
		o.enableMetrics = true
	}
}

// WithDebug enables debug mode.
// Debug mode allows the caching system to log debug information.
func WithDebug() optionsFunc {
	return func(o *Options) {
		o.enableDebug = true
	}
}

// WithDebugf sets a custom debug log function in the optionsuration.
// This function is responsible for logging debug messages.
func WithDebugf(fn func(format string, v ...any)) optionsFunc {
	return func(o *Options) {
		o.debugf = fn
	}
}
