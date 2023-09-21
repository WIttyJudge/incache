package cache

import "sync/atomic"

type Metrics struct {
	// Shows how many times items were inserted into cache.
	insertions uint64

	// Shows how many times items were successfully retrived by key.
	hits uint64

	// Shows how many times items weren't retrived by key.
	misses uint64

	// Shows how many items were released from the cache.
	evictions uint64
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

// Get collected insertions.
func (m *Metrics) Insertions() uint64 {
	return atomic.LoadUint64(&m.insertions)
}

// Get collected hits.
func (m *Metrics) Hits() uint64 {
	return atomic.LoadUint64(&m.hits)
}

// Get collected misses.
func (m *Metrics) Misses() uint64 {
	return atomic.LoadUint64(&m.misses)
}

// Get collected evictions.
func (m *Metrics) Evictions() uint64 {
	return atomic.LoadUint64(&m.evictions)
}

func (m *Metrics) incrInsertions() {
	atomic.AddUint64(&m.insertions, 1)
}

func (m *Metrics) incrHits() {
	atomic.AddUint64(&m.hits, 1)
}

func (m *Metrics) incrMisses() {
	atomic.AddUint64(&m.misses, 1)
}

func (m *Metrics) incrEvictions() {
	atomic.AddUint64(&m.evictions, 1)
}
