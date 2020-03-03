package cache

import (
	"sync"
	"time"
)

type SyncLRU struct {
	lru LRUCache
	mu  *sync.Mutex
}

func NewSyncLRU(capacity int, ttl time.Duration) LRUCache {
	return &SyncLRU{
		lru: NewLRU(capacity, ttl),
		mu:  &sync.Mutex{},
	}
}

func (slru *SyncLRU) SetClock(clock Clock) {
	slru.mu.Lock()
	defer slru.mu.Unlock()
	slru.lru.SetClock(clock)
}

func (slru *SyncLRU) Exists(key string) bool {
	slru.mu.Lock()
	defer slru.mu.Unlock()
	return slru.lru.Exists(key)
}

func (slru *SyncLRU) Set(key string, value interface{}) {
	slru.mu.Lock()
	defer slru.mu.Unlock()
	slru.lru.Set(key, value)
}

func (slru *SyncLRU) Delete(key string) {
	slru.mu.Lock()
	defer slru.mu.Unlock()
	slru.lru.Delete(key)
}

func (slru *SyncLRU) Get(key string) (interface{}, bool) {
	slru.mu.Lock()
	defer slru.mu.Unlock()
	return slru.lru.Get(key)
}

func (slru *SyncLRU) TTL(key string) (time.Duration, bool) {
	slru.mu.Lock()
	defer slru.mu.Unlock()
	return slru.lru.TTL(key)
}

func (slru *SyncLRU) Expired() int {
	slru.mu.Lock()
	defer slru.mu.Unlock()
	return slru.lru.Expired()
}

func (slru *SyncLRU) Evicted() int {
	slru.mu.Lock()
	defer slru.mu.Unlock()
	return slru.lru.Evicted()
}

func (slru *SyncLRU) UpdateTTL(update bool) {
	slru.mu.Lock()
	defer slru.mu.Unlock()
	slru.lru.UpdateTTL(update)
}

func (slru *SyncLRU) OnEvict(callback func(key string)) {
	slru.mu.Lock()
	defer slru.mu.Unlock()
	slru.lru.OnEvict(callback)
}

func (slru *SyncLRU) OnExpire(callback func(key string)) {
	slru.mu.Lock()
	defer slru.mu.Unlock()
	slru.lru.OnExpire(callback)
}
