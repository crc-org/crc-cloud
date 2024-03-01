package fx

import "sort"

func Sorted[T any](it Iterator[T], less func(a, b T) bool) Iterator[T] {
	s := ToSlice(it)
	sort.Slice(s, func(i, j int) bool { return less(s[i], s[j]) })
	return IterSlice(s)
}
