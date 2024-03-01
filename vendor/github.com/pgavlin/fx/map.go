package fx

type fmap[T, U any] struct {
	it Iterator[T]
	fn func(v T) (U, bool)
	v  U
}

func (f *fmap[T, U]) Value() U {
	return f.v
}

func (f *fmap[T, U]) Next() bool {
	for {
		if !f.it.Next() {
			var v U
			f.v = v
			return false
		}
		if v, ok := f.fn(f.it.Value()); ok {
			f.v = v
			return true
		}
	}
}

func OfType[T, U any](it Iterator[T]) Iterator[U] {
	return FMap[T, U](it, func(v T) (U, bool) {
		u, ok := ((interface{})(v)).(U)
		return u, ok
	})
}

func FMap[T, U any](it Iterator[T], fn func(v T) (U, bool)) Iterator[U] {
	return &fmap[T, U]{it: it, fn: fn}
}

func Map[T, U any](it Iterator[T], fn func(v T) U) Iterator[U] {
	return FMap[T, U](it, func(v T) (U, bool) { return fn(v), true })
}
