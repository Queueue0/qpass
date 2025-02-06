package structures

import "fmt"

type node[T any] struct {
	v T
	next *node[T]
}

type Queue[T any] struct {
	head *node[T]
	tail *node[T]
}

func (q *Queue[T]) Enqueue(e T) {
	n := &node[T]{e, nil}
	if q.head == nil {
		q.head = n
		q.tail = q.head
	} else {
		q.tail.next = n
		q.tail = n
	}
}

func (q *Queue[T]) Dequeue() T {
	t := q.head.v
	q.head = q.head.next
	return t
}

func (q *Queue[T]) HasNext() bool {
	return q.head != nil
}

func (q *Queue[T]) String() string {
	out := "Queue: ["
	n := q.head
	for n != nil {
		out += fmt.Sprintf("%v", n.v)
		if n.next != nil  {
			out += ", "
		}
		n = n.next
	}
	out += "]"
	return out
}
