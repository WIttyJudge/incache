package cache

import "time"

type Options struct {
	ttl           time.Duration
	enableMetrics bool
}

type optionsFunc func(*Options)

func defaultOptions() Options {
	return Options{
		ttl:           5 * time.Minute,
		enableMetrics: false,
	}
}

// Sets default TTL for all items that whould be stored in the cache.
func WithTTL(ttl time.Duration) optionsFunc {
	return func(opts *Options) {
		opts.ttl = ttl
	}
}

// Enables the collection of metrics that run throughout the cache work.
func WithMetrics() optionsFunc {
	return func(opts *Options) {
		opts.enableMetrics = true
	}
}
