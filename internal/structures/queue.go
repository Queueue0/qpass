package structures

type Queue[T any] struct {
	q []T
}

func (q *Queue[T]) Push(e T) {
	q.q = append(q.q, e)
}

func (q *Queue[T]) Pop() T {
	t := q.q[0]
	q.q = q.q[1:]
	return t
}

func (q *Queue[T]) HasNext() bool {
	return len(q.q) > 0
}
