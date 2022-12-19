package collections

type Set[T comparable] struct {
	Map map[T]bool
}

func NewSet[T comparable]() Set[T] {
	return Set[T]{Map: make(map[T]bool)}
}

func (s *Set[T]) Contains(k T) bool {
	if s.Map == nil {
		s.Map = make(map[T]bool)
	}
	_, ok := s.Map[k]
	return ok
}

func (s *Set[T]) Remove(k T) {
	if s.Map == nil {
		s.Map = make(map[T]bool)
	}
	delete(s.Map, k)
}

func (s *Set[T]) Set(k T) {
	if s.Map == nil {
		s.Map = make(map[T]bool)
	}
	s.Map[k] = true
}

func (s *Set[T]) Keys() []T {
	if s.Map == nil {
		s.Map = make(map[T]bool)
	}
	keys := make([]T, 0, len(s.Map))
	for k := range s.Map {
		keys = append(keys, k)
	}
	return keys
}

func (s *Set[T]) Len() int {
	return len(s.Map)
}
