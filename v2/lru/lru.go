package lru

import (
	"time"
)

type Cache interface {
	Capacity() int
	Exists(key string) bool
	Set(key string, value interface{})
	Delete(key string) bool
	Get(key string) (interface{}, bool)
	TTL(key string) (time.Duration, bool)
	Destroy()
}

type item struct {
	data     interface{}
	expireAt time.Time
}
