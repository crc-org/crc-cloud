package fx

type take[T any] struct {
	n  int
	it Iterator[T]
	v  T
}

func (t *take[T]) Value() T {
	return t.v
}

func (t *take[T]) Next() bool {
	if t.n <= 0 {
		return false
	}
	if !t.it.Next() {
		var v T
		t.v = v
		return false
	}
	t.v, t.n = t.it.Value(), t.n-1
	return true
}

func Take[T any](it Iterator[T], n int) Iterator[T] {
	return &take[T]{n: n, it: it}
}
