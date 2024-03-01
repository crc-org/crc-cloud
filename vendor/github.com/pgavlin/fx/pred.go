package fx

func Any[T any](it Iterator[T], pred func(v T) bool) bool {
	for it.Next() {
		if pred(it.Value()) {
			return true

		}
	}
	return false
}

func All[T any](it Iterator[T], pred func(v T) bool) bool {
	for it.Next() {
		if !pred(it.Value()) {
			return false
		}
	}
	return true
}
