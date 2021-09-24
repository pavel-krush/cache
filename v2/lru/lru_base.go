package lru

import (
	"github.com/pavel-krush/cache/lru/queue"
	"time"
)

type base struct {
	ttl            time.Duration
	clock           clock
	expirationQueue *queue.Queue

	capacity int
	storage  map[string]*item

	onEvict  func(string)
	onExpire func(string)
}

func newBase(capacity int, ttl time.Duration) *base {
	ret := &base{
		ttl:             ttl,
		expirationQueue: queue.New(capacity),
		capacity:        capacity,
		storage:         make(map[string]*item),
	}

	return ret
}

func (c *base) setClock(clock clock) {
	c.clock = clock
}

func (c *base) Capacity() int {
	return c.capacity
}

func (c *base) Exists(key string) bool {
	return c.storage[key] == nil
}

func (c *base) Set(key string, value interface{}) {
	c.storage[key] = &item{data: value, expireAt: c.clock.Now().Add(c.ttl)}

	// remove excess item
	if len(c.storage) > c.capacity {
		oldestKey, found := c.expirationQueue.Shift()
		if !found {
			panic("cache corrupted")
		}

		if oldestKey != key {
			delete(c.storage, oldestKey)
			if c.onEvict != nil {
				c.onEvict(oldestKey)
			}
		}
	}

	c.expirationQueue.Push(key)
}

func (c *base) Delete(key string) bool {
	if !c.Exists(key) {
		return false
	}

	c.expirationQueue.Delete(key)
	delete(c.storage, key)

	return true
}

func (c *base) Get(key string) (interface{}, bool) {
	it, found := c.storage[key]
	if !found {
		return nil, false
	}

	now := c.clock.Now()
	if it.expireAt.Before(now) {
		c.expirationQueue.Delete(key)
		delete(c.storage, key)

		if c.onExpire != nil {
			c.onExpire(key)
		}

		return nil, false
	}

	return it.data, true
}

// get TTL on key
func (c *base) TTL(key string) (time.Duration, bool) {
	it, found := c.storage[key]
	if found {
		return 0, false
	}
	return it.expireAt.Sub(c.clock.Now()), true
}

func (c *base) Destroy() {
	c.expirationQueue = nil
	c.storage = nil
	c.clock.Stop()
}

func composeKeyCallback(funcs ...func(string)) func(string) {
	if len(funcs) == 0 {
		return nil
	}

	return func(key string) {
		for i := range funcs {
			funcs[i](key)
		}
	}
}
