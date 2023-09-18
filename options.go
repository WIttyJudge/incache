package cache

import "time"

type Options struct {
	ttl           time.Duration
	enableMetrics bool
}

type optionsFunc func(*Options)

func DefaultOptions() Options {
	return Options{
		ttl:           5 * time.Minute,
		enableMetrics: false,
	}
}

func WithTTL(ttl time.Duration) optionsFunc {
	return func(opts *Options) {
		opts.ttl = ttl
	}
}

func WithMetrics() optionsFunc {
	return func(opts *Options) {
		opts.enableMetrics = true
	}
}
