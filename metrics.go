package cache

import "sync/atomic"

type Metrics struct {
	// Show how many times items were inserted into cache
	insertions uint64

	// Show how many times items were successfully retrived by key
	hits uint64

	// Show how many times items weren't retrived by key
	misses uint64
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

// Get insertions
func (m *Metrics) Insertions() uint64 {
	return atomic.LoadUint64(&m.insertions)
}

// Get hits
func (m *Metrics) Hits() uint64 {
	return atomic.LoadUint64(&m.hits)
}

// get misses
func (m *Metrics) Misses() uint64 {
	return atomic.LoadUint64(&m.misses)
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
