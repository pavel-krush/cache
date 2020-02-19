package cache

import (
	"time"
)

type LRUCache interface {
	Exists(key string) bool
	Set(key string, value interface{})
	Delete(key string)
	Get(key string) (interface{}, bool)
	TTL(key string) (time.Duration, bool)
	Expired() int
	Evicted() int
	UpdateTTL(update bool)

	OnEvict(func(key string))
	OnExpire(func(key string))
}

type LRU struct {
	ttl            time.Duration
	clock          Clock
	expirationList *l

	evicted int
	expired int

	capacity  int
	updateTTL bool
	storage   map[string]*Item

	onEvict func (key string)
	onExpire func (key string)
}

type Item struct {
	data     interface{}
	expireAt time.Time
}

func NewLRU(capacity int, ttl time.Duration) LRUCache {
	return &LRU{
		ttl:            ttl,
		clock:          ClockSimple,
		expirationList: newList(capacity),
		capacity:       capacity,
		storage:        make(map[string]*Item),
	}
}

func (lru *LRU) Exists(key string) bool {
	if _, ok := lru.storage[key]; ok {
		return true
	}
	return false
}

func (lru *LRU) Set(key string, value interface{}) {
	lru.expire()

	item := &Item{data: value, expireAt: time.Now().Add(lru.ttl)}
	lru.storage[key] = item

	// remove excess item
	if len(lru.storage) >= lru.capacity {
		lru.evict()
	}

	lru.expirationList.insert(key)
}

func (lru *LRU) Delete(key string) {
	lru.expirationList.delete(key)
	delete(lru.storage, key)
}

func (lru *LRU) Get(key string) (interface{}, bool) {
	lru.expire()

	item, found := lru.storage[key]
	if !found {
		return nil, false
	}

	if lru.updateTTL {
		item.expireAt = time.Now().Add(lru.ttl)
		lru.expirationList.moveToFront(key)
	}

	return item.data, true
}

// get TTL on key
func (lru *LRU) TTL(key string) (time.Duration, bool) {
	if lru.clock == ClockNone {
		return 1<<63 - 1, true
	}

	item, found := lru.storage[key]
	if !found {
		return 0, false
	}
	return item.expireAt.Sub(lru.clock.Now()), true
}

func (lru *LRU) Expired() int {
	return lru.expired
}

func (lru *LRU) Evicted() int {
	return lru.evicted
}

// remove the oldest element
func (lru *LRU) evict() {
	key, evicted := lru.expirationList.pop()
	if evicted {
		lru.evicted++
		delete(lru.storage, key)

		if lru.onEvict != nil {
			lru.onEvict(key)
		}
	}
}

// remove all expired elements
func (lru *LRU) expire() {
	if lru.clock == ClockNone {
		return
	}
	for {
		oldestKey, peeked := lru.expirationList.peek()
		if !peeked {
			break
		}
		now := lru.clock.Now()
		item := lru.storage[oldestKey]
		// stop at first not expired element
		if !item.expireAt.Before(now) {
			break
		}
		lru.expirationList.pop()
		delete(lru.storage, oldestKey)
		lru.expired++

		if lru.onExpire != nil {
			lru.onExpire(oldestKey)
		}
	}
}

func (lru *LRU) UpdateTTL(update bool) {
	lru.updateTTL = update
}

func (lru *LRU) OnEvict(callback func (key string)) {
	lru.onEvict = callback
}

func (lru *LRU) OnExpire(callback func (key string)) {
	lru.onExpire = callback
}