package xsf

import (
	"container/heap"
	"fmt"
	"sort"
	"sync/atomic"
	"testing"
)

func Test_priQueue(t *testing.T) {
	pq := PriorityQueue{
		data: []*Item{
			newItem("banana", 3),
			newItem("apple", 2),
			newItem("pear", 4),
		},
	}

	heap.Init(&pq)

	item := &Item{value: func() atomic.Value { var value atomic.Value; value.Store("orange"); return value }(), priority: 1}

	heap.Push(&pq, item)
	pq.update(item, item.value.Load().(string), 5)

	fmt.Printf("x addr:%+v\n", pq.load(0))
	fmt.Printf("x addr:%+v\n", pq.load(2))

	fmt.Println()

	pq.delete(item)
	pq.delete(item)

	fmt.Printf("addr:%+v\n", pq.load(0))

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("%.2d:%s ", item.priority, item.value)
	}
}
func Test_items(t *testing.T) {
	data := Items([]*Item{
		newItem("banana", 3),
		newItem("apple", 2),
		newItem("pear", 4),
	})
	sort.Sort(data)
	for _, item := range data {
		t.Log(*item)
	}
	var i int
	var j *Item
	for i, j = range data {
		if j.value.Load().(string) == "apple" {
			break
		}
	}
	data.Delete(i)
	for _, item := range data {
		t.Log(*item)
	}
}
func Test_Queue(t *testing.T) {
	data := []*Item{
		newItem("banana", 3),
		newItem("apple", 2),
		newItem("pear", 4),
	}
	queue := newQueue(data)

	t.Log(queue.load(0))
	queue.update(nil, "apple", 9)
	t.Log('x',queue.load(0))
	t.Log('x',queue.load(2))
}
