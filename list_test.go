package cache

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func queueDump(q *l) {
	var elements []string
	if q.tail == 0 {
		elements = append(elements, "queue: <empty>")
	} else {
		elements = append(elements, fmt.Sprintf("queue: <head %d, tail %d>", q.head, q.tail))

		ptr := q.head
		for {
			item := q.list[ptr]
			elements = append(elements, fmt.Sprintf("  % 2d: % 2d %s % 2d", ptr, item.left, item.key, item.right))

			ptr = q.list[ptr].right
			if ptr == q.head {
				break
			}
		}
	}
	elements = append(elements, "")
	_, _ = fmt.Fprintf(os.Stderr, strings.Join(elements, "\n"))
}

func queueKeys(q *l) []string {
	var ret []string
	if q.tail == 0 {
		return ret
	}

	ptr := q.head
	for {
		item := q.list[ptr]
		ret = append(ret, item.key)
		ptr = item.right
		if ptr == q.head {
			break
		}
	}
	return ret
}

type testcaseOp string

const (
	opAdd    testcaseOp = "add"
	opDelete            = "delete"
)

type queueTestCase struct {
	op       testcaseOp
	what     string
	elements []string
}

func TestQueue(t *testing.T) {
	q := newList(10)

	testcases := []queueTestCase{
		{op: opAdd, what: "one", elements: []string{"one"}},
		{op: opAdd, what: "two", elements: []string{"two", "one"}},
		{op: opAdd, what: "three", elements: []string{"three", "two", "one"}},
		{op: opDelete, what: "two", elements: []string{"three", "one"}},
		{op: opAdd, what: "two", elements: []string{"two", "three", "one"}},
		{op: opDelete, what: "one", elements: []string{"two", "three"}},
		{op: opAdd, what: "three", elements: []string{"three", "two"}},
		{op: opDelete, what: "three", elements: []string{"two"}},
		{op: opDelete, what: "two", elements: []string{}},
	}

	for _, testcase := range testcases {
		if testcase.op == opAdd {
			q.insert(testcase.what)
		} else if testcase.op == opDelete {
			q.delete(testcase.what)
		} else {
			panic("unknown testcase operator")
		}
		actualKeys := queueKeys(q)
		if len(actualKeys) != len(testcase.elements) {
			t.Errorf("%v %v", actualKeys, testcase.elements)
			continue
		}

		for i := range actualKeys {
			if actualKeys[i] != testcase.elements[i] {
				t.Errorf("%v %v", actualKeys, testcase.elements)
				break
			}
		}
	}
}
