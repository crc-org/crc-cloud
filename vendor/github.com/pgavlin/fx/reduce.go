package fx

func Reduce[T, U any](it Iterator[T], init U, fn func(acc U, v T) U) U {
	for it.Next() {
		init = fn(init, it.Value())
	}
	return init
}
