package cache

import (
	"math/rand"
	"testing"
	"time"
)

func makeRandomWord(l int) string {
	var alphabet = "qazwsxedcrfvtgbyhnujmikolpQAZWSXEDCRFVTGBYHNUJMIKOLP"
	ret := make([]byte, l)
	for i := 0; i < l; i++ {
		ret[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(ret)
}

const accessKeysSize = 1000000

func BenchmarkMapNoExpiration(b *testing.B) {
	cache := make(map[string]int)
	var keys []string
	var accessKeys = make([]string, accessKeysSize)
	for i := 0; i < 10000; i++ {
		key := makeRandomWord(10)
		keys = append(keys, key)
		cache[key] = rand.Intn(1024)
	}

	for i := 0; i < accessKeysSize; i++ {
		accessKeys[i] = keys[rand.Intn(len(keys))]
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cache[accessKeys[i%accessKeysSize]]
	}
}

func benchmarkLru(b *testing.B, cache LRUCache) {
	var keys []string
	var accessKeys = make([]string, accessKeysSize)
	for i := 0; i < 10000; i++ {
		key := makeRandomWord(10)
		keys = append(keys, key)
		cache.Set(key, rand.Intn(1024))
	}

	for i := 0; i < accessKeysSize; i++ {
		accessKeys = append(accessKeys, keys[rand.Intn(len(keys))])
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(accessKeys[i%accessKeysSize])
	}
}

func BenchmarkLRUNoExpiration(b *testing.B) {
	cache := NewLRU(10000, time.Hour)
	benchmarkLru(b, cache)
}

func BenchmarkSyncLRUNoExpiration(b *testing.B) {
	cache := NewSyncLRU(10000, time.Hour)
	benchmarkLru(b, cache)
}
