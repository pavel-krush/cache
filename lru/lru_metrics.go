package lru

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// lruWithMetrics is a wrapper for cache that exports prometheus metrics
type lruWithMetrics struct {
	parent Cache

	capacityMetric prometheus.Gauge
	hitsMetric     prometheus.Counter
	missesMetric   prometheus.Counter
	evictedMetric  prometheus.Counter
	expiredMetric  prometheus.Counter
}

func newWithMetrics(
	parent Cache,
	namespace string,
	subsystem string,
	constLabels prometheus.Labels,
) *lruWithMetrics {
	capacity := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   namespace,
		Subsystem:   subsystem,
		Name:        "cache_capacity",
		Help:        "Maximum number of items in cache",
		ConstLabels: constLabels,
	})
	capacity.Set(float64(parent.Capacity()))

	hits := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   namespace,
		Subsystem:   subsystem,
		Name:        "cache_hits_total",
		Help:        "Total amount of cache hits",
		ConstLabels: constLabels,
	})

	misses := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   namespace,
		Subsystem:   subsystem,
		Name:        "cache_misses_total",
		Help:        "Total amount of cache misses",
		ConstLabels: constLabels,
	})

	evicted := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   namespace,
		Subsystem:   subsystem,
		Name:        "cache_evicted_total",
		Help:        "Total amount of keys evicted by cache overflow",
		ConstLabels: constLabels,
	})

	expired := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   namespace,
		Subsystem:   subsystem,
		Name:        "cache_expired_total",
		Help:        "Total amount of expired keys",
		ConstLabels: constLabels,
	})

	prometheus.MustRegister(capacity)
	prometheus.MustRegister(hits)
	prometheus.MustRegister(misses)
	prometheus.MustRegister(evicted)
	prometheus.MustRegister(expired)

	return &lruWithMetrics{
		parent: parent,

		capacityMetric: capacity,
		hitsMetric:     hits,
		missesMetric:   misses,
		evictedMetric:  evicted,
		expiredMetric:  expired,
	}
}

func (c *lruWithMetrics) Capacity() int {
	return c.parent.Capacity()
}

func (c *lruWithMetrics) Exists(key string) bool {
	exists := c.parent.Exists(key)
	if exists {
		c.hitsMetric.Inc()
	} else {
		c.missesMetric.Inc()
	}

	return exists
}

func (c *lruWithMetrics) Set(key string, value interface{}) {
	c.parent.Set(key, value)
}

func (c *lruWithMetrics) Delete(key string) bool {
	deleted := c.parent.Delete(key)
	if deleted {
		c.hitsMetric.Inc()
	} else {
		c.missesMetric.Inc()
	}

	return deleted
}

func (c *lruWithMetrics) Get(key string) (interface{}, bool) {
	ret, found := c.parent.Get(key)
	if found {
		c.hitsMetric.Inc()
	} else {
		c.missesMetric.Inc()
	}

	return ret, found
}

func (c *lruWithMetrics) TTL(key string) (time.Duration, bool) {
	ttl, found := c.parent.TTL(key)
	if found {
		c.hitsMetric.Inc()
	} else {
		c.missesMetric.Inc()
	}

	return ttl, found
}

func (c *lruWithMetrics) Destroy() {
	prometheus.Unregister(c.capacityMetric)
	prometheus.Unregister(c.hitsMetric)
	prometheus.Unregister(c.missesMetric)
	prometheus.Unregister(c.evictedMetric)
	prometheus.Unregister(c.expiredMetric)

	c.parent.Destroy()
}

func (c *lruWithMetrics) onEvict(string) {
	c.evictedMetric.Inc()
}

func (c *lruWithMetrics) onExpire(string) {
	c.expiredMetric.Inc()
}
