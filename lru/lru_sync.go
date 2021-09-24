package lru

import (
	"sync"
	"time"
)

// lruWithSync is a wrapper for cache that allows concurrent access to the cache
type lruWithSync struct {
	parent Cache

	sync.Mutex
}

func newWithSync(parent Cache) *lruWithSync {
	return &lruWithSync{parent: parent}
}

func (c *lruWithSync) Capacity() int {
	c.Lock()
	ret := c.parent.Capacity()
	c.Unlock()

	return ret
}

func (c *lruWithSync) Exists(key string) bool {
	c.Lock()
	ret := c.parent.Exists(key)
	c.Unlock()

	return ret
}

func (c *lruWithSync) Set(key string, value interface{}) {
	c.Lock()
	c.parent.Set(key, value)
	c.Unlock()
}

func (c *lruWithSync) Delete(key string) bool {
	c.Lock()
	ret := c.parent.Delete(key)
	c.Unlock()

	return ret
}

func (c *lruWithSync) Get(key string) (interface{}, bool) {
	c.Lock()
	val, ok := c.parent.Get(key)
	c.Unlock()

	return val, ok
}

func (c *lruWithSync) TTL(key string) (time.Duration, bool) {
	c.Lock()
	val, ok := c.parent.TTL(key)
	c.Unlock()

	return val, ok
}

func (c *lruWithSync) Destroy() {
	c.parent.Destroy()
}
