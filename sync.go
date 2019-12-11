package cache

import (
	"sync"
	"time"
)

type SyncLRU struct {
	LRU
	mu *sync.Mutex
}

func NewSyncLRU(capacity int, ttl time.Duration, clock Clock) *SyncLRU {
	return &SyncLRU{
		LRU: *NewLRU(capacity, ttl, clock),
		mu:  &sync.Mutex{},
	}
}

func (slru *SyncLRU) Exists(key string) bool {
	slru.mu.Lock()
	defer slru.mu.Unlock()
	return slru.LRU.Exists(key)
}

func (slru *SyncLRU) Set(key string, value interface{}) {
	slru.mu.Lock()
	defer slru.mu.Unlock()
	slru.LRU.Set(key, value)
}

func (slru *SyncLRU) Delete(key string) {
	slru.mu.Lock()
	defer slru.mu.Unlock()
	slru.LRU.Delete(key)
}

func (slru *SyncLRU) Get(key string) (interface{}, bool) {
	slru.mu.Lock()
	defer slru.mu.Unlock()
	return slru.LRU.Get(key)
}

func (slru *SyncLRU) TTL(key string) (time.Duration, bool) {
	slru.mu.Lock()
	defer slru.mu.Unlock()
	return slru.LRU.TTL(key)
}
