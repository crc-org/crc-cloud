package fx

func Filter[T any](it Iterator[T], fn func(v T) bool) Iterator[T] {
	return FMap[T, T](it, func(v T) (T, bool) { return v, fn(v) })
}
