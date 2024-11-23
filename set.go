package armap

type setValue struct{}

type Set[K comparable] struct {
	m *Map[K, setValue]
}

func (s *Set[K]) Len() int {
	return s.m.Len()
}

func (s *Set[K]) Add(key K) bool {
	_, ok := s.m.Set(key, setValue{})
	return ok
}

func (s *Set[K]) Contains(key K) bool {
	_, ok := s.m.Get(key)
	return ok
}

func (s *Set[K]) Delete(key K) bool {
	_, ok := s.m.Delete(key)
	return ok
}

func (s *Set[K]) Scan(iter func(K) bool) {
	s.m.Scan(func(key K, value setValue) bool {
		return iter(key)
	})
}

func (s *Set[K]) Clear() {
	s.m.Clear()
}

func (s *Set[K]) Release() {
	s.m.Release()
}

func NewSet[K comparable](arena Arena, funcs ...OptionFunc) *Set[K] {
	a := NewTypeArena[Set[K]](arena)
	return a.NativeNewValue(func(s *Set[K]) {
		s.m = NewMap[K, setValue](arena, funcs...)
	})
}
