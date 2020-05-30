package util

import (
	"container/heap"
)

type Item struct {
	value interface{}
	index int
}

func (i *Item) Value() interface{} {
	return i.value
}

type innerPriorityQueue struct {
	items    []*Item
	lessFunc func(interface{}, interface{}) bool
}

type PriorityQueue struct {
	ipq innerPriorityQueue
}

func NewPriorityQueue(cmpFunc func(interface{}, interface{}) bool) *PriorityQueue {
	var pq PriorityQueue
	pq.ipq.items = make([]*Item, 0)
	pq.ipq.lessFunc = cmpFunc
	heap.Init(&pq.ipq)
	return &pq
}

func (pq *PriorityQueue) Len() int {
	return pq.ipq.Len()
}

func (pq *PriorityQueue) Push(x interface{}) *Item {
	i := new(Item)
	i.value = x
	heap.Push(&pq.ipq, i)
	heap.Fix(&pq.ipq, i.index)
	return i
}

func (pq *PriorityQueue) Pop() interface{} {
	item := heap.Pop(&pq.ipq).(*Item)
	return item.value
}

func (pq *PriorityQueue) Peek() interface{} {
	return pq.ipq.items[0].value
}

func (pq *PriorityQueue) Update(item *Item, value interface{}) {
	item.value = value
	heap.Fix(&pq.ipq, item.index)
}

func (pq *innerPriorityQueue) Len() int {
	return len(pq.items)
}

func (pq *innerPriorityQueue) Less(i, j int) bool {
	return pq.lessFunc(pq.items[i].value, pq.items[j].value)
}

func (pq *innerPriorityQueue) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].index = i
	pq.items[j].index = j
}

func (pq *innerPriorityQueue) Push(x interface{}) {
	n := len(pq.items)
	i := x.(*Item)
	i.index = n
	pq.items = append(pq.items, i)
}

func (pq *innerPriorityQueue) Pop() interface{} {
	old := pq.items
	n := len(old)
	item := old[n-1]
	item.index = -1
	pq.items = old[0 : n-1]
	return item
}
