package fx

type List[T any] interface {
	Len() int
	At(i int) T
}

type list[T any] []T

func (l list[T]) Len() int {
	return len(l)
}

func (l list[T]) At(i int) T {
	return l[i]
}

func AsList[T any](ts []T) List[T] {
	return list[T](ts)
}
