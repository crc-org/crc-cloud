package fx

type concat[T any] struct {
	iters Iterator[Iterator[T]]
	it    Iterator[T]
	v     T
}

func (e *concat[T]) Value() T {
	return e.v
}

func (e *concat[T]) Next() bool {
	for {
		if e.it != nil {
			if e.it.Next() {
				e.v = e.it.Value()
				return true
			}
		}

		if !e.iters.Next() {
			var v T
			e.v = v
			return false
		}

		e.it = e.iters.Value()
	}
	return false
}

func Concat[T any](iters ...Iterator[T]) Iterator[T] {
	return ConcatMany(IterSlice(iters))
}

func ConcatMany[T any](iters Iterator[Iterator[T]]) Iterator[T] {
	return &concat[T]{iters: iters}
}
