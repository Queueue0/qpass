package structures

type node[T any] struct {
	v T
	next *node[T]
}

type Queue[T any] struct {
	head *node[T]
	tail *node[T]
}

func (q *Queue[T]) Enqueue(e T) {
	if q.head == nil {
		q.head = &node[T]{e, nil}
		q.tail = q.head
	} else {
		n := &node[T]{e, nil}
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
