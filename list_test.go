package cache

import (
	"testing"
)

func queueKeys(q *l) []string {
	var ret []string
	if q.tail == 0 {
		return ret
	}

	ptr := 0
	for {
		item := q.list[ptr]
		ret = append(ret, item.key)
		ptr = item.right
		if ptr == 0 {
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

func testQueue(t *testing.T) {
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
