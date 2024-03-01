package fx

type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(v T) {
	s[v] = struct{}{}
}

func (s Set[T]) In(v T) bool {
	_, ok := s[v]
	return ok
}
