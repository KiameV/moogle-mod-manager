package collections

type Set interface {
	Add(key string)
	Contains(key string) bool
	Remove(key string) bool
	Size() int
}

type set struct {
	data map[string]bool
}

func NewSet() Set {
	return &set{make(map[string]bool)}
}

func (s *set) Add(key string) {
	s.data[key] = true
}

func (s *set) Contains(key string) (exists bool) {
	_, exists = s.data[key]
	return
}

func (s *set) Remove(key string) (exists bool) {
	_, exists = s.data[key]
	return
}

func (s *set) Size() int {
	return len(s.data)
}
