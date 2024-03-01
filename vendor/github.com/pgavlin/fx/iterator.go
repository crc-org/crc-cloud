package fx

import "reflect"

type Iterator[T any] interface {
	Value() T
	Next() bool
}

type Pair[T, U any] struct {
	Fst T
	Snd U
}

func (p Pair[T, U]) Unpack() (T, U) {
	return p.Fst, p.Snd
}

func NewPair[T, U any](fst T, snd U) Pair[T, U] {
	return Pair[T, U]{Fst: fst, Snd: snd}
}

func ToSlice[T any](it Iterator[T]) []T {
	var s []T
	for it.Next() {
		s = append(s, it.Value())
	}
	return s
}

func TrySlice[T any](it Iterator[Result[T]]) ([]T, error) {
	var s []T
	for it.Next() {
		v, err := it.Value().Unpack()
		if err != nil {
			return nil, err
		}
		s = append(s, v)
	}
	return s, nil
}

func ToSet[T comparable](it Iterator[T]) Set[T] {
	s := Set[T]{}
	for it.Next() {
		s.Add(it.Value())
	}
	return s
}

func TrySet[T comparable](it Iterator[Result[T]]) (Set[T], error) {
	s := Set[T]{}
	for it.Next() {
		v, err := it.Value().Unpack()
		if err != nil {
			return Set[T]{}, err
		}
		s.Add(v)
	}
	return s, nil
}

func ToMap[K comparable, V any](it Iterator[Pair[K, V]]) map[K]V {
	m := map[K]V{}
	for it.Next() {
		key, value := it.Value().Unpack()
		m[key] = value
	}
	return m
}

func TryMap[K comparable, V any](it Iterator[Result[Pair[K, V]]]) (map[K]V, error) {
	m := map[K]V{}
	for it.Next() {
		v, err := it.Value().Unpack()
		if err != nil {
			return nil, err
		}
		key, value := v.Unpack()
		m[key] = value
	}
	return m, nil
}

func IterSlice[T any](ts []T) Iterator[T] {
	return IterList(AsList(ts))
}

func IterList[T any](ts List[T]) Iterator[T] {
	return &iterator[T]{ts: ts}
}

func IterMap[K comparable, V any](m map[K]V) Iterator[Pair[K, V]] {
	return &mapIterator[K, V]{it: reflect.ValueOf(m).MapRange()}
}

func IterSet[T comparable](s Set[T]) Iterator[T] {
	return Map(IterMap(map[T]struct{}(s)), func(p Pair[T, struct{}]) T {
		return p.Fst
	})
}

type mapIterator[K comparable, V any] struct {
	it *reflect.MapIter
	v  Pair[K, V]
}

func (i *mapIterator[K, V]) Value() Pair[K, V] {
	return i.v
}

func (i *mapIterator[K, V]) Next() bool {
	if !i.it.Next() {
		i.v = Pair[K, V]{}
		return false
	}
	i.v = Pair[K, V]{
		Fst: i.it.Key().Interface().(K),
		Snd: i.it.Value().Interface().(V),
	}
	return true
}

type only[T any] struct {
	v    T
	done bool
}

func (o *only[T]) Value() T {
	return o.v
}

func (o *only[T]) Next() bool {
	if !o.done {
		o.done = true
		return true
	}
	return false
}

func Only[T any](v T) Iterator[T] {
	return &only[T]{v: v}
}

type empty[T any] struct{}

func (empty[T]) Value() (v T) {
	return
}

func (empty[T]) Next() bool {
	return false
}

func Empty[T any]() Iterator[T] {
	return empty[T]{}
}

type iterator[T any] struct {
	idx int
	v   T
	ts  List[T]
}

func (i *iterator[T]) Value() (v T) {
	return i.v
}

func (i *iterator[T]) Next() bool {
	if i.idx >= i.ts.Len() {
		var v T
		i.v = v
		return false
	}
	i.v = i.ts.At(i.idx)
	i.idx++
	return true
}
