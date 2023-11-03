package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/wittyjudge/incache"
)

func main() {
	cache := incache.New(incache.WithMetrics())
	performCacheOperations(cache)

	registerAndExposeMetrics(cache)
	startMetricsServer()
}

func performCacheOperations(cache *incache.Cache) {
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	cache.Get("key1")
	cache.Get("key4")
	cache.Delete("key2")
}

func registerAndExposeMetrics(cache *incache.Cache) {
	// Get the cache metrics
	metrics := cache.Metrics()

	// Register cache metrics with Prometheus
	prometheus.MustRegister(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "incache_items_inserted_total",
			Help: "Number of items inserted",
		},
		func() float64 {
			return float64(metrics.Insertions())
		}))

	prometheus.MustRegister(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "incache_items_hitted_total",
			Help: "Number of items hitted",
		},
		func() float64 {
			return float64(metrics.Hits())
		}))

	prometheus.MustRegister(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "incache_items_missed_total",
			Help: "Number of items missed",
		},
		func() float64 {
			return float64(metrics.Misses())
		}))

	prometheus.MustRegister(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "incache_items_evicted_total",
			Help: "Number of items evicted",
		},
		func() float64 {
			return float64(metrics.Evictions())
		}))

	prometheus.MustRegister(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "incache_items_count_current",
			Help: "Number of items currently stored in cache",
		},
		func() float64 {
			return float64(cache.Len())
		}))
}

func startMetricsServer() {
	// Expose the registered metrics via HTTP
	http.Handle("/metrics", promhttp.Handler())

	// Start the server
	http.ListenAndServe(":8080", nil)
}
