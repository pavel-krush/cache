package cache

import (
	"testing"
)

func TestQueue(t *testing.T) {
	q := NewQueue(10)
	q.Unshift("one")
	q.Unshift("two")
	q.Unshift("three")
	q.Unshift("four")
	q.Unshift("five")
	q._dump()
}
