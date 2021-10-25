package lru

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Builder struct {
	optCapacity        *optionCapacity
	optTTL             *optionTTL
	optSync            *optionSync
	optMetrics         *optionMetrics
	optDiscreteClock   *optionDiscreteClock
	optEvictCallbacks  []*optionEvictCallback
	optExpireCallbacks []*optionExpireCallback
}

func New() Builder {
	return Builder{}
}

func NewFromConfig(cfg *Config) Builder {
	cfg = cfg.withDefaults()

	ret := New().WithCapacity(cfg.Capacity).
		WithTTL(cfg.TTL)

	if cfg.Concurrent {
		ret = ret.WithSync()
	}

	if cfg.Metrics != nil {
		ret = ret.WithMetrics(cfg.Metrics.Namespace, cfg.Metrics.Subsystem, cfg.Metrics.Labels)
	}

	if cfg.Clock != nil {
		if cfg.Clock.Discrete != nil {
			ret = ret.WithDiscreteClock(cfg.Clock.Discrete.UpdateInterval)
		}
	}

	return ret
}

func (b Builder) WithCapacity(capacity int) Builder {
	if b.optCapacity != nil {
		panic("duplicated WithCapacity()")
	}

	b.optCapacity = &optionCapacity{capacity}
	return b
}

func (b Builder) WithTTL(ttl time.Duration) Builder {
	if b.optTTL != nil {
		panic("duplicated WithTTL()")
	}

	b.optTTL = &optionTTL{ttl}
	return b
}

func (b Builder) WithSync() Builder {
	if b.optSync != nil {
		panic("duplicated WithSync()")
	}

	b.optSync = &optionSync{}
	return b
}

func (b Builder) WithMetrics(namespace string, subsystem string, constLabels prometheus.Labels) Builder {
	if b.optMetrics != nil {
		panic("duplicated WithMetrics()")
	}

	b.optMetrics = &optionMetrics{namespace, subsystem, constLabels}

	return b
}

func (b Builder) WithDiscreteClock(updateInterval time.Duration) Builder {
	if b.optDiscreteClock != nil {
		panic("duplicated WithDiscreteClock()")
	}

	b.optDiscreteClock = &optionDiscreteClock{updateInterval}
	return b
}

func (b Builder) WithEvictCallback(cb func(string)) Builder {
	b.optEvictCallbacks = append(b.optEvictCallbacks, &optionEvictCallback{cb})
	return b
}

func (b Builder) WithExpireCallback(cb func(string)) Builder {
	b.optExpireCallbacks = append(b.optExpireCallbacks, &optionExpireCallback{cb})
	return b
}

func (b Builder) Build() Cache {
	// capacity is mandatory
	if b.optCapacity == nil || b.optCapacity.capacity <= 0 {
		panic("LRU cache capacity must be greater than zero")
	}

	// nil ttl is same as ttl = 0
	if b.optTTL == nil {
		b.optTTL = &optionTTL{0}
	}

	if b.optTTL.ttl < 0 {
		panic("LRU cache TTL must be greater or equal to zero")
	}

	var (
		onEvictCallbacks  []func(string)
		onExpireCallbacks []func(string)
	)

	baseCache := newBase(b.optCapacity.capacity, b.optTTL.ttl)

	if b.optTTL.ttl == 0 {
		baseCache.setClock(newClockNone())
	} else {
		if b.optDiscreteClock == nil {
			baseCache.setClock(newClockPrecise())
		} else {
			baseCache.setClock(newClockDiscrete(b.optDiscreteClock.updateInterval))
		}
	}

	var ret Cache = baseCache

	if b.optMetrics != nil {
		withMetrics := newWithMetrics(ret,
			b.optMetrics.namespace, b.optMetrics.subsystem, b.optMetrics.constLabels)

		onEvictCallbacks = append(onEvictCallbacks, withMetrics.onEvict)
		onExpireCallbacks = append(onExpireCallbacks, withMetrics.onExpire)

		ret = withMetrics
	}

	if b.optSync != nil {
		ret = newWithSync(ret)
	}

	for i := range b.optEvictCallbacks {
		onEvictCallbacks = append(onEvictCallbacks, b.optEvictCallbacks[i].cb)
	}

	for i := range b.optExpireCallbacks {
		onExpireCallbacks = append(onExpireCallbacks, b.optExpireCallbacks[i].cb)
	}

	baseCache.onEvict = composeKeyCallback(onEvictCallbacks...)
	baseCache.onExpire = composeKeyCallback(onExpireCallbacks...)

	return ret
}
