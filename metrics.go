package incache

import "sync/atomic"

// Metrics stores cache statistics
type metrics struct {
	// Shows how many times items were inserted into cache.
	insertions uint64

	// Shows how many times items were successfully retrieved by key.
	hits uint64

	// Shows how many times items weren't retrieved by key.
	misses uint64

	// Shows how many items were released from the cache.
	evictions uint64
}

func newMetrics() *metrics {
	return &metrics{}
}

// Get collected insertions.
func (m *metrics) Insertions() uint64 {
	return atomic.LoadUint64(&m.insertions)
}

// Get collected hits.
func (m *metrics) Hits() uint64 {
	return atomic.LoadUint64(&m.hits)
}

// Get collected misses.
func (m *metrics) Misses() uint64 {
	return atomic.LoadUint64(&m.misses)
}

// Get collected evictions.
func (m *metrics) Evictions() uint64 {
	return atomic.LoadUint64(&m.evictions)
}

func (m *metrics) incrInsertions() {
	atomic.AddUint64(&m.insertions, 1)
}

func (m *metrics) incrHits() {
	atomic.AddUint64(&m.hits, 1)
}

func (m *metrics) incrMisses() {
	atomic.AddUint64(&m.misses, 1)
}

func (m *metrics) incrEvictions() {
	atomic.AddUint64(&m.evictions, 1)
}
