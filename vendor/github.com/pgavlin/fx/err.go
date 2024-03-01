package fx

type Result[T any] struct {
	v   T
	err error
}

func (r Result[T]) Unpack() (T, error) {
	return r.v, r.err
}

func OK[T any](v T) Result[T] {
	return Result[T]{v: v}
}

func Err[T any](e error) Result[T] {
	return Result[T]{err: e}
}

func Try[T any](v T, e error) Result[T] {
	if e != nil {
		return Err[T](e)
	}
	return OK(v)
}
