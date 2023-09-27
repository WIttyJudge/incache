package incache

import (
	"testing"
)

func BenchmarkSet(b *testing.B) {
	cache := New()
	defer cache.Close()

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cache.Set("key", "value")
		}
	})
}

func BenchmarkGet(b *testing.B) {
	cache := New()
	defer cache.Close()

	cache.Set("key0", "value")

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = cache.Get("key0")
		}
	})
}
