package incache

import (
	"fmt"
	"testing"
	"time"
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
			cache.Get("key0")
		}
	})
}

func BenchmarkDeleteExpired(b *testing.B) {
	cache := New(WithTTL(1 * time.Millisecond))
	defer cache.Close()

	for i := 0; i < 1000000; i++ {
		key := fmt.Sprint(i)
		cache.Set(key, i)
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cache.DeleteExpired()
		}
	})
}
