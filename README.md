# GO LRU TTL Cache 
Simple LRU cache implementation. All keys have same TTL and cache has limited capacity.

## Basic usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/pavel-krush/cache/v2/lru"
)

func main() {
	cache := lru.New().WithCapacity(10000).WithTTL(time.Second * 60).Build()

	cache.Set("key", "value")
	value, found := cache.Get("key")
	fmt.Printf("%s = %s, %t", "key", value, found)
}
```

## Build options
There are several build options that can produce cache with different behaviour:
- WithCapacity(capacity int). Mandatory. Sets the capacity of the cache;
- WithTTL(ttl time.Duration). Optional. Set the keys TTL. All keys have same TTL;
- WithMetrics(namespace string, subsystem string, constLabels []string). Optional. Creates cache with metrics  
  The following metrics will be registered:  
  - namespace_subsystem_cache_capacity{constLabels} - Gauge: capacity of the cache.;
  - namespace_subsystem_cache_hits_total{constLabels} - Counter: amount of cache hits;
  - namespace_subsystem_cache_misses_total{constLabels} - Counter: amount of cache misses;
  - namespace_subsystem_cache_evicted_total{constLabels} - Counter: amount of evicted keys;
  - namespace_subsystem_cache_expired_total{constLabels} - Counter: amount of expired keys.  
  
  Metrics are registered on cache creation and de-registered when cache is destroyed via `.Destroy()`.
- WithSync(). Optional. Creates a concurrent cache.
- WithDiscreteClock(time.Duration). Optional. Creates a cache with less precise clock.  
  This option allows to increase performance of `.Get()`.
- WithEvictCallback(func(string)). Optional. Adds an eviction hook(see below);
- WithExpireCallback(func(string)). Optional. Adds an expiration hook(see below).

## Cache config
It's possible to build cache from config structure:

```go
cfg := lru.Config{
	Capacity:   314,
	TTL:        42 * time.Second,
	Concurrent: true,
	Metrics: &lru.MetricsConfig{
		Enabled:   true,
		Namespace: "namespace",
		Subsystem: "subsystem",
		Labels: prometheus.Labels{
			"name": "cache_name",
		},
	},
	Clock: &lru.ClockConfig{
		Discrete: &lru.ClockConfigDiscrete{UpdateInterval: 500 * time.Millisecond},
	},
}

cache := lru.NewFromConfig(cfg).Build()
```

Only `Capacity` and `TTL` fields are mandatory.

Default values:
- `Concurrent` false
- `Metrics` { Enabled: false }
- `Clock` { Simple: {} }

## Hooks

### Eviction
Eviction hook is called when a new key is pushed into the full cache. In this case the oldest key will be evicted from the cache. 
It's okay to create several eviction callbacks.

### Expiration
Expire hook is called when key is removed from the cache and its TTL has passed.
It's not guaranteed that this hook will be called just in time when key is expired.
It's okay to create several expiration callbacks.

LRU Cache Interface
---------

```go
type Cache interface {
    Capacity() int
    Exists(key string) bool
    Set(key string, value interface{})
    Delete(key string) bool
    Get(key string) (interface{}, bool)
    TTL(key string) (time.Duration, bool)
    Destroy()
}
```
