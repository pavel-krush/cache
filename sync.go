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
	slru.lru.SetClock(clock)
	slru.mu.Unlock()
}

func (slru *SyncLRU) Exists(key string) bool {
	slru.mu.Lock()
	ret := slru.lru.Exists(key)
	slru.mu.Unlock()
	return ret
}

func (slru *SyncLRU) Set(key string, value interface{}) {
	slru.mu.Lock()
	slru.lru.Set(key, value)
	slru.mu.Unlock()
}

func (slru *SyncLRU) Delete(key string) {
	slru.mu.Lock()
	slru.lru.Delete(key)
	slru.mu.Unlock()
}

func (slru *SyncLRU) Get(key string) (interface{}, bool) {
	slru.mu.Lock()
	val, ok := slru.lru.Get(key)
	slru.mu.Unlock()
	return val, ok
}

func (slru *SyncLRU) TTL(key string) (time.Duration, bool) {
	slru.mu.Lock()
	val, ok := slru.lru.TTL(key)
	slru.mu.Unlock()
	return val, ok
}

func (slru *SyncLRU) GetStats() Stats {
	slru.mu.Lock()
	ret := slru.lru.GetStats()
	slru.mu.Unlock()
	return ret
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
