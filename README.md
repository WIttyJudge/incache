# incache

[![Go Report Card](https://goreportcard.com/badge/github.com/wittyjudge/incache)](https://goreportcard.com/report/github.com/wittyjudge/incache)
[![MIT License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/WIttyJudge/incache/blob/main/LICENSE)

Simple thread-safe time-based caching library for Go.

The library includes:

- Automatic removal of expired data, which can be disabled easily;
- Optional collection of metrics;

## Installation

```bash
go get github.com/wittyjudge/incache
```

## Usage

### Simple Initialization

```go
package main

import (
	"fmt"
	"time"

	"github.com/wittyjudge/incache"
)

func main() {
	cache := incache.New()

	cache.Set("key1", "value1")
	cache.SetWithTTL("key2", "value2", 1*time.Minute)

	value := cache.Get("key1")
	fmt.Println(value)

	cache.Delete("key2")
}
```

### Custom Initialization

Note that by default, a new cache instance runs with default config.
You can find its default options [here](https://github.com/WIttyJudge/incache/blob/main/config.go#L16).
However, passing config options into the `incache.New()` allows you to set desired
behavior.

```go
package main

import (
	"time"

	"github.com/wittyjudge/incache"
)

func main() {
	// incache.WithTTL() sets TTL for all items that would be stored in the cache.
	// TTL <= 0 means that the item won't have expiration time at all.
	//
	// incache.WithCleanupInterval() sets Interval between removing expired items.
	// If the interval is less than or equal to 0, no automatic clearing
	// is performed.
	//
	// incache.WithMetrics() enables the collection of metrics that run throughout the cache work.
	cache := incache.New(
		incache.WithTTL(5*time.Minute),
		incache.WithCleanupInterval(5*time.Minute),
		incache.WithMetrics(),
	)

	cache.Set("key1", "value1")
}
```

## Benchmarks

```
go test -bench=. -benchmem -benchtime=4s

goos: linux
goarch: amd64
pkg: github.com/wittyjudge/incache
cpu: Intel(R) Core(TM) i5-8350U CPU @ 1.70GHz
BenchmarkSet-8          28206325               168.2 ns/op             0 B/op          0 allocs/op
BenchmarkGet-8          90718786                52.57 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/wittyjudge/incache   2.331s

```

## Development Roadmap

- [x] Add collection of metrics;
- [ ] Tests and benchmarks;
- [ ] Examples of usage in documentation;
- [ ] Subscribe to events like eviction and insertion;

## License

`incache` source code is available under the MIT License.
