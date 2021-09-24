package queue

import (
	"fmt"
	"testing"
)

func Benchmark_queue(b *testing.B) {
	size := b.N

	q := New(size)
	keys := make([]string, size)

	for i := 0; i < size; i++ {
		keys[i] = fmt.Sprintf("key-%d", i)
		q.Push(keys[i])
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		q.Push(keys[i%size])
	}
}
