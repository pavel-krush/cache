package cache

import (
	"fmt"
	"os"
	"strings"
)

type Queue struct {
	capacity int
	keys     map[string]int
	list     []QueueItem
	head     int // pointer to current head. it always valid except the case when tail == 0. this means queue is empty
	tail     int // pointer to the next free element
}

type QueueItem struct {
	left, right int
	key string
}

func NewQueue(capacity int) *Queue {
	ret := &Queue{
		capacity: capacity,
		keys:     make(map[string]int),
		list:     make([]QueueItem, capacity),
	}
	return ret
}

func (q *Queue) Size() int {
	return q.capacity
}

func (q *Queue) Exists(key string) bool {
	_, e := q.keys[key]
	return e
}

func (q *Queue) Unshift(key string) {
	if q.Exists(key) {
		q.moveToFront(key)
		return
	}

	if q.tail == q.capacity {
		panic("queue full")
	}

	item := QueueItem{
		left: q.list[q.head].left,
		right: q.head,
		key: key,
	}

	// uuh
	q.list[q.tail] = item
	q.list[q.list[q.head].left].right = q.tail
	q.list[q.head].left = q.tail
	q.head = q.tail
	q.keys[key] = q.tail
	q.tail++
}

func (q *Queue) moveToFront(key string) {

}

func (q *Queue) _dump() {
	var elements []string
	if q.head == q.tail {
		elements = append(elements, "queue: <empty>")
	} else {
		elements = append(elements, fmt.Sprintf("queue: <head %d, tail %d>", q.head, q.tail))
	}
	ptr := q.head
	for {
		item := q.list[ptr]
		elements = append(elements, fmt.Sprintf("  % 2d: % 2d %s % 2d", ptr, item.left, item.key, item.right))

		ptr = q.list[ptr].right
		if ptr == q.head {
			break
		}
	}
	elements = append(elements, "")
	fmt.Fprintf(os.Stderr, strings.Join(elements, "\n"))
}
