GO LRU TTL Cache
================
This is a simple LRU TTL cache implementation. Cache has limited capacity and keys expire after certain time. 

Usage
-----

```go
import "github.com/pavel-krush/cache"

cache := cache.NewLRU(100000, time.Second * 60)

cache.Set("key", "value")
value, found := cache.Get("key")
```

Interface
---------

```go

func NewLRU(capacity int, ttl time.Duration) LRUCache
func NewSyncLRU(capacity int, ttl time.Duration) LRUCache

type LRUCache interface {
	Exists(key string) bool // check whether key exists in cache
	Set(key string, value interface{}) // set key-value pair
	Delete(key string) // delete key from cache
	Get(key string) (interface{}, bool) // get value from cache
	TTL(key string) (time.Duration, bool) // get TTL on key
	Expired() int // get total count of expired elements
	Evicted() int // get total count of evicted elements
	UpdateTTL(update bool) // update or not element's ttl on Get()

        // eviction and expiration callbacks
	OnEvict(func(key string))
	OnExpire(func(key string))
}
```

Thread safety
-------------

`SyncLRU` is a thread safe version of `LRU` with exact the same interface.

Benchmarks
----------

Benchmark example for my MacBook Pro (13-inch, 2019), 2,4 GHz Quad-Core Intel Core i5:
```text
BenchmarkMapNoExpiration-8       	30545750	        36.7 ns/op
BenchmarkLRUNoExpiration-8       	17854893	        67.1 ns/op
BenchmarkSyncLRUNoExpiration-8   	11962401	        99.1 ns/op
```

Clock
-----
There are three types of clock:
- ClockNone - Used for non-expiration cache
- ClockSimple - High precision clock. It uses time.Now() on each Get()
- DiscreteClock - Not as precise as ClockSimple, but significantly faster. Refreshes current time once in 500ms. Default Clock.