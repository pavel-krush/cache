package cache

type l struct {
	capacity int
	keys     map[string]int
	list     []listItem
	// pointer to current head. it always valid except the case when tail == 0. this means list is empty
	head int
	// pointer to the next free element
	tail int
}

type listItem struct {
	left, right int
	key         string
}

func newList(capacity int) *l {
	ret := &l{
		capacity: capacity,
		keys:     make(map[string]int),
		list:     make([]listItem, capacity),
	}
	return ret
}

func (q *l) size() int {
	return q.tail
}

func (q *l) insert(key string) {
	if _, ok := q.keys[key]; ok {
		q.moveToFront(key)
		return
	}

	if q.tail == q.capacity {
		panic("list full")
	}

	item := listItem{
		left:  q.list[q.head].left,
		right: q.head,
		key:   key,
	}

	// insert new element into the end of the list
	q.list[q.tail] = item

	// advance pointers to make new element to be first in the list
	q.list[q.list[q.head].left].right = q.tail
	q.list[q.head].left = q.tail

	// advance head pointer
	q.head = q.tail

	q.keys[key] = q.tail
	q.tail++
}

func (q *l) delete(key string) {
	index, ok := q.keys[key]
	if !ok {
		return
	}

	delete(q.keys, key)
	q.tail--

	// advance head if current head is being removed
	if q.head == index {
		q.head = q.list[q.head].right
	}

	// detach current element from list
	q.list[q.list[index].left].right = q.list[index].right
	q.list[q.list[index].right].left = q.list[index].left

	if index == q.tail {
		return
	}

	// move tail element to freed position
	q.list[q.list[q.tail].left].right = index
	q.list[q.list[q.tail].right].left = index

	q.list[index] = q.list[q.tail]
	q.keys[q.list[index].key] = index

	if q.head == q.tail {
		q.head = index
	}

	// q.list[q.tail] = listItem{}
}

func (q *l) moveToFront(key string) {
	q.delete(key)
	q.insert(key)
}

func (q *l) pop() (string, bool) {
	if q.tail == 0 {
		return "", false
	}
	key := q.list[q.list[q.head].left].key
	q.delete(key)
	return key, true
}

func (q *l) peek() (string, bool) {
	if q.tail == 0 {
		return "", false
	}
	return q.list[q.list[q.head].left].key, true
}
