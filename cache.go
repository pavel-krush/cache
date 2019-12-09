package cache

import "time"

type Item struct {
	data interface {}
	expireAt time.Time
}

type Cache struct {
	ttl time.Duration
	clock *Clock
}
//
//func (c *Cache) Exists(key string) bool {
//
//}
//
//func (c *Cache) Delete(key string) {
//
//}
//
//func (c *Cache) Get(key string) (interface{}, bool) {
//
//}
//
//func (c *Cache) Set(key string, value interface{}) {
//
//}
//
//func (c *Cache) TTL(key string) (time.Duration, bool) {
//
//}
