package lru

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func key(i int) string {
	return fmt.Sprintf("key-%d", i)
}

func value(i int) string {
	return fmt.Sprintf("value-%d", i)
}

func randomWord(l int) string {
	var alphabet = "qazwsxedcrfvtgbyhnujmikolpQAZWSXEDCRFVTGBYHNUJMIKOLP"
	ret := make([]byte, l)
	for i := 0; i < l; i++ {
		ret[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(ret)
}

func Test_LRU_base_basic(t *testing.T) {
	capacity := 10
	c := New().WithCapacity(capacity).Build()

	for i := 0; i < capacity; i++ {
		k := key(i)
		v := value(i)

		c.Set(k, v)
	}

	for i := 0; i < capacity; i++ {
		k := key(i)
		v := value(i)

		val, found := c.Get(k)
		if !found {
			t.Errorf("key \"%s\" not found", key(i))
			continue
		}

		valStr := val.(string)

		if valStr != v {
			t.Errorf("expected \"%s\" = \"%s\", got \"%s\"", k, v, valStr)
		}
	}
}

func Test_LRU_base_eviction(t *testing.T) {
	capacity := 5

	c := New().WithCapacity(capacity).Build()

	for i := 0; i < capacity; i++ {
		c.Set(key(i), value(i))
	}

	c.Set(key(capacity), value(capacity))

	_, found := c.Get(key(0))
	if found {
		t.Errorf("expected key \"%s\" evicted", key(0))
	}
}

func Test_LRU_base_expiration(t *testing.T) {
	capacity := 10
	ttl := time.Millisecond * 50

	c := New().WithCapacity(capacity).WithTTL(ttl).Build()

	for i := 0; i < capacity; i++ {
		c.Set(key(i), value(i))
	}

	for i := 0; i < capacity; i++ {
		if _, found := c.Get(key(i)); !found {
			t.Errorf("expected key \"%s\" in cache", key(i))
		}
	}

	time.Sleep(ttl)

	for i := 0; i < capacity; i++ {
		if _, found := c.Get(key(i)); found {
			t.Errorf("expected key \"%s\" expired", key(i))
		}
	}
}

const accessKeysSize = 1000000

func BenchmarkMapNoExpiration(b *testing.B) {
	cache := make(map[string]int)
	var keys []string
	var accessKeys = make([]string, accessKeysSize)
	for i := 0; i < 10000; i++ {
		k := randomWord(10)
		keys = append(keys, k)
		cache[k] = rand.Intn(1024)
	}

	for i := 0; i < accessKeysSize; i++ {
		accessKeys[i] = keys[rand.Intn(len(keys))]
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cache[accessKeys[i%accessKeysSize]]
	}
}

func benchmarkLru(b *testing.B, cache Cache) {
	var keys []string
	var accessKeys = make([]string, accessKeysSize)
	for i := 0; i < 10000; i++ {
		k := randomWord(10)
		keys = append(keys, k)
		cache.Set(k, rand.Intn(1024))
	}

	for i := 0; i < accessKeysSize; i++ {
		accessKeys[i] = keys[rand.Intn(len(keys))]
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(accessKeys[i%accessKeysSize])
	}
}

func BenchmarkLRUNoExpiration(b *testing.B) {
	cache := New().
		WithCapacity(10000).
		WithTTL(time.Hour).
		Build()
	benchmarkLru(b, cache)
}

func BenchmarkSyncLRUNoExpiration(b *testing.B) {
	cache := New().WithCapacity(10000).WithSync().WithTTL(time.Hour).Build()
	benchmarkLru(b, cache)
}
