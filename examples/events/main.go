package main

import (
	"fmt"
	"time"

	"github.com/wittyjudge/incache"
)

func main() {
	cache := incache.New()

	// It's important to call the Close method because it ensures
	// that all concurrently running events are waited for until the end of execution.
	defer cache.Close()

	cache.OnInsertion(func(key string, value interface{}) {
		time.Sleep(300 * time.Millisecond)
		fmt.Printf("Insertion event was triggered: key: %s, value: %v\n", key, value)
	})

	cache.OnEviction(func(key string, value interface{}) {
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("Eviction event was triggered: key: %s, value: %v\n", key, value)
	})

	fmt.Println("Performing Set operation")
	cache.Set("key1", 123)
	cache.Set("key2", 1235)
	cache.Set("key3", 100)

	fmt.Println("Performing Delete operation")
	cache.Delete("key1")
	cache.Delete("key3")
}
