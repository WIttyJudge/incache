# incache

[![GoDoc](https://godoc.org/github.com/wittyjudge/incache?status.png)](https://godoc.org/github.com/wittyjudge/incache)
[![codecov](https://codecov.io/gh/WIttyJudge/incache/graph/badge.svg)](https://codecov.io/gh/WIttyJudge/incache)
[![Release](https://img.shields.io/github/v/release/wittyjudge/incache)](https://github.com/WIttyJudge/incache/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/wittyjudge/incache)](https://goreportcard.com/report/github.com/wittyjudge/incache)

Simple thread-safe time-based caching library for Go.

## Features

- Automatic removal of expired data, which can be disabled easily;
- Collection of metrics;
- Debug mode;
- Event handlers (insertion and eviction)

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
	// Initialize cache instance with default config.
	cache := incache.New()

	// Set a new value.
	cache.Set("key1", "value1")
	// Set a new value with 1 minute expiration time.
	cache.SetWithTTL("key2", "value2", 1*time.Minute)

	// Get the value for the key 'key1'
	value := cache.Get("key1")
	fmt.Println(value)

	// Delete the value for the key 'key2'
	cache.Delete("key2")
}
```

### Configuration Options

Note that by default, a new cache instance runs with default config.
However, passing config options into the `incache.New()` allows you to set desired
behavior.

#### TTL

Defines the default TTL for all items that would be stored in
the cache. TTL <= 0 means that the item won't have expiration time at all.
The default value is 5 minutes.

Example:

```go
cache := incache.New(incache.WithTTL(30 * time.Minute))
```

#### CleanupInterval

Defines the interval between removing expired items.
If the interval is less than or equal to 0, no automatic clearing is performed.
The default value is 5 minutes.

Example:

```go
cache := incache.New(incache.WithCleanupInterval(30 * time.Minute))
```

#### WithMetrics

Enables metrics collection.
Allows for the continuous collection of metrics during the cache's operation.
The default value is false.

Example:

```go
cache := incache.New(incache.WithMetrics())
```

#### WithDebug 

Enables debug mode.
Allowing the logging of debug information to the stdout.
The default value is false.

Example:

```go
cache := incache.New(incache.WithDebug())
```

#### WithDebugf

Defines a custom debug log function.
This function is responsible for logging debug messages.
The default value uses the internal Go log package with the '[incache]' prefix set at the beginning.

Example:

```go
myLogFunc := log.New(os.Stdout, "[myprefix]", 0).Printf
cache := incache.New(incache.WithDebugf(myLogFunc))
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

- [x] Cache metrics (at least hits, insertions, misses, evictions rate);
- [x] Tests and benchmarks;
- [x] Subscribe to events like eviction and insertion;
- [ ] Change code that evicts expired items to use priority queue;

## License

`incache` source code is available under the MIT License.
