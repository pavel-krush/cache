package cache

import (
	"sync"
	"sync/atomic"
	"time"
)

type Clock interface {
	Now() time.Time
	Reset()
}

// Fake clock. Used for non-expiration cache
type ClockNone struct{}

func (cn ClockNone) Now() time.Time {
	return time.Time{}
}
func (cn ClockNone) Reset() {}

func NewClockNone() Clock {
	return &ClockNone{}
}

// Simple clock. Used for precise expiration
type ClockSimple struct{}

func (cs ClockSimple) Now() time.Time {
	return time.Now()
}
func (cs ClockSimple) Reset() {}

func NewClockSimple() Clock {
	return &ClockSimple{}
}

// Optimized clock. Not precise as SimpleClick but significantly faster
type ClockDiscrete struct {
	mu    *sync.Mutex
	value *atomic.Value
}

func (cd ClockDiscrete) Now() time.Time {
	now := cd.value.Load().(time.Time)
	return now
}

func (cd ClockDiscrete) Reset() {}

func (cd ClockDiscrete) Refresh() {
	cd.value.Store(time.Now())
}

func NewClockDiscrete(updateTime time.Duration) Clock {
	if updateTime == 0 {
		updateTime = time.Millisecond * 500
	}
	ret := ClockDiscrete{
		mu:    &sync.Mutex{},
		value: &atomic.Value{},
	}
	ret.Refresh()
	go func() {
		ticker := time.NewTicker(updateTime)
		for {
			<-ticker.C
			ret.Refresh()
		}
	}()
	return ret
}
