package cache

import (
	"sync"
	"time"
)

type Cache struct {
	mu             *sync.Mutex
	ttl            time.Duration
	clock          Clock
	expirationList *l

	evicted int
	expired int

	capacity int
	storage  map[string]Item
}

type Item struct {
	data     interface{}
	expireAt time.Time
}

func NewCache(capacity int, ttl time.Duration, clock Clock) *Cache {
	return &Cache{
		mu:             &sync.Mutex{},
		ttl:            ttl,
		clock:          clock,
		expirationList: newList(capacity),
		capacity:       capacity,
		storage:        make(map[string]Item),
	}
}

func (c *Cache) Exists(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.storage[key]; ok {
		return true
	}
	return false
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.expire()

	item := Item{data: value, expireAt: time.Now().Add(c.ttl)}
	c.storage[key] = item

	// remove excess item
	if len(c.storage) >= c.capacity {
		c.evict()
	}

	c.expirationList.insert(key)
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.expirationList.delete(key)
	delete(c.storage, key)
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.expire()

	item, found := c.storage[key]
	if !found {
		return nil, false
	}

	return item.data, true
}

// get TTL on key
func (c *Cache) TTL(key string) (time.Duration, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.clock == ClockNone {
		return 1<<63 - 1, true
	}

	item, found := c.storage[key]
	if !found {
		return 0, false
	}
	return item.expireAt.Sub(c.clock.Now()), true
}

// remove the oldest element
func (c *Cache) evict() {
	key, evicted := c.expirationList.pop()
	if evicted {
		c.evicted++
		delete(c.storage, key)
	}
}

// remove all expired elements
func (c *Cache) expire() {
	if c.clock == ClockNone {
		return
	}
	for {
		oldestKey, peeked := c.expirationList.peek()
		if !peeked {
			break
		}
		now := c.clock.Now()
		item := c.storage[oldestKey]
		// stop at first not expired element
		if !item.expireAt.Before(now) {
			break
		}
		c.expirationList.pop()
		delete(c.storage, oldestKey)
		c.expired++
	}
}
