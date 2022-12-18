package collections

type Set[T comparable] struct {
	m map[T]bool
}

func NewSet[T comparable]() Set[T] {
	return Set[T]{m: make(map[T]bool)}
}

func (s *Set[T]) Contains(k T) bool {
	if s.m == nil {
		s.m = make(map[T]bool)
	}
	_, ok := s.m[k]
	return ok
}

func (s *Set[T]) Remove(k T) {
	if s.m == nil {
		s.m = make(map[T]bool)
	}
	delete(s.m, k)
}

func (s *Set[T]) Set(k T) {
	if s.m == nil {
		s.m = make(map[T]bool)
	}
	s.m[k] = true
}

func (s *Set[T]) Keys() []T {
	if s.m == nil {
		s.m = make(map[T]bool)
	}
	keys := make([]T, 0, len(s.m))
	for k := range s.m {
		keys = append(keys, k)
	}
	return keys
}
