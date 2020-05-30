package util

type (
	Queue struct {
		length     int
		head, tail *qnode
	}

	qnode struct {
		entry interface{}
		next  *qnode
	}
)

func (n *qnode) Value() interface{} {
	return n.entry
}

func NewQueue() *Queue {
	return &Queue{length: 0, head: nil, tail: nil}
}

func (q *Queue) Len() int {
	return q.length
}

func (q *Queue) Push(entry interface{}) {
	node := &qnode{entry: entry}
	if q.length == 0 {
		q.head = node
		q.tail = node
	} else {
		q.tail.next = node
		q.tail = node
	}
	q.length++
}

func (q *Queue) Pop() interface{} {
	if q.length == 0 {
		return nil
	}
	node := q.head
	if q.length == 1 {
		q.head = nil
		q.tail = nil
	} else {
		q.head = q.head.next
	}
	q.length--
	return node.entry
}

func (q *Queue) Peek() interface{} {
	if q.length == 0 {
		return nil
	}
	return q.head.entry
}
