package incache

import "sync/atomic"

type metrics interface {
	Insertions() uint64
	Hits() uint64
	Misses() uint64
	Evictions() uint64

	reset()

	incrementInsertions()
	incrementHits()
	incrementMisses()
	incrementEvictions()
}

// Metrics stores cache statistics
type realMetrics struct {
	// Shows how many times items were inserted into cache.
	insertions uint64

	// Shows how many times items were successfully retrieved by key.
	hits uint64

	// Shows how many times items weren't retrieved by key.
	misses uint64

	// Shows how many items were released from the cache.
	evictions uint64
}

func newRealMetrics() *realMetrics {
	return &realMetrics{}
}

// Get collected insertions.
func (m *realMetrics) Insertions() uint64 {
	return atomic.LoadUint64(&m.insertions)
}

// Get collected hits.
func (m *realMetrics) Hits() uint64 {
	return atomic.LoadUint64(&m.hits)
}

// Get collected misses.
func (m *realMetrics) Misses() uint64 {
	return atomic.LoadUint64(&m.misses)
}

// Get collected evictions.
func (m *realMetrics) Evictions() uint64 {
	return atomic.LoadUint64(&m.evictions)
}

func (m *realMetrics) reset() {
	m.insertions = 0
	m.hits = 0
	m.misses = 0
	m.evictions = 0
}

func (m *realMetrics) incrementInsertions() {
	m.insertions += 1
}

func (m *realMetrics) incrementHits() {
	m.hits += 1
}

func (m *realMetrics) incrementMisses() {
	m.misses += 1
}

func (m *realMetrics) incrementEvictions() {
	m.evictions += 1
}

// Dummy metrics implementation that is used if metrics is disabled.
type noMetrics struct{}

func newNoMetrics() *noMetrics { return &noMetrics{} }

func (m *noMetrics) Insertions() uint64 { return 0 }
func (m *noMetrics) Hits() uint64       { return 0 }
func (m *noMetrics) Misses() uint64     { return 0 }
func (m *noMetrics) Evictions() uint64  { return 0 }

func (m *noMetrics) reset() {}

func (m *noMetrics) incrementInsertions() {}
func (m *noMetrics) incrementHits()       {}
func (m *noMetrics) incrementMisses()     {}
func (m *noMetrics) incrementEvictions()  {}
