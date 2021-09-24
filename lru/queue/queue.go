package queue

import (
	"fmt"
	"sort"
)

// Queue is a simple limited size queue implementation
type Queue struct {
	keys map[string]int
	list []queueItem
	free []int // free items
	head int // index of first head element
}

type queueItem struct {
	left  int // left points to the element that is older than current
	right int // right points to the element that is newer than current
	key   string
}

func New(capacity int) *Queue {
	ret := &Queue{
		keys: make(map[string]int),
		list: make([]queueItem, capacity),
		free: make([]int, capacity),
	}

	// mark all elements as free
	for i := 0; i < capacity; i++ {
		ret.free[i] = i
	}

	return ret
}

// Push puts a new element into the end of the queue
func (q *Queue) Push(key string) {
	if _, ok := q.keys[key]; ok {
		q.MoveToEnd(key)
		return
	}

	freeLen := len(q.free)

	isFull := freeLen == 0
	isEmpty := freeLen == cap(q.free)

	if isFull {
		panic("queue full")
	}

	index := q.free[freeLen - 1]
	q.free = q.free[:freeLen - 1]

	q.keys[key] = index
	q.list[index].key = key

	if isEmpty {
		q.head = index
		q.list[index].left = index
		q.list[index].right = index
	} else {
		// new element should become last in the queue
		// it's left must point to the previous last element
		// it's right must point to the head
		q.list[index].left = q.list[q.head].left
		q.list[index].right = q.head

		// remember previous last element
		prevLast := q.list[q.head].left

		// move head's left to new element
		q.list[q.head].left = index

		// move previous last element's right to new element
		q.list[prevLast].right = index
	}
}

// Shift extracts the first element from the queue
func (q *Queue) Shift() (string, bool) {
	if len(q.free) == cap(q.free) {
		return "", false
	}

	key := q.list[q.head].key
	q.Delete(key)

	return key, true
}

// Delete deletes given element from the queue
func (q *Queue) Delete(key string) {
	index, ok := q.keys[key]
	if !ok {
		return
	}

	delete(q.keys, key)

	// detach current element from list
	q.list[q.list[index].left].right = q.list[index].right
	q.list[q.list[index].right].left = q.list[index].left

	// when head is being deleted, move the head to the next element
	if index == q.head {
		q.head = q.list[index].right
	}

	q.free = append(q.free, index)
}

// Peek returns the first element of the queue
func (q *Queue) Peek() (string, bool) {
	isEmpty := len(q.free) == cap(q.free)
	if isEmpty {
		return "", false
	}

	return q.list[q.head].key, true
}

// MoveToEnd makes given element to be the last element in the queue
func (q *Queue) MoveToEnd(key string) {
	q.Delete(key)
	q.Push(key)
}

func (q *Queue) DebugPrint() {
	fmt.Printf("head: %d\n", q.head)

	fmt.Printf("free map: (%p) [", q.free)
	for i := 0; i < len(q.free); i++ {
		fmt.Printf("%d", q.free[i])
		if i < len(q.free) - 1 {
			fmt.Printf(", ")
		}
	}
	fmt.Printf("]\n")

	fmt.Printf("idx  left     key right  free\n")
	for i := 0; i < len(q.list); i++ {
		free := false
		for freeIdx := 0; freeIdx < len(q.free); freeIdx++ {
			if q.free[freeIdx] == i {
				free = true
			}
		}

		fmt.Printf("%3d  %3d %7s %3d     %t\n",
			i,
			q.list[i].left,
			"\"" + q.list[i].key + "\"",
			q.list[i].right,
			free,
		)
	}

	keys := make([]string, 0, len(q.keys))
	for k := range q.keys {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return q.keys[keys[i]] < q.keys[keys[j]]
	})

	fmt.Printf("keys:\n")
	for i := range keys {
		fmt.Printf("%3s: %d\n", keys[i], q.keys[keys[i]])
	}

	fmt.Printf("\n")
}
