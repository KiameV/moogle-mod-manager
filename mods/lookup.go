package mods

import "fmt"

type (
	lookupID string
	mod      interface {
		ID() ModID
		Kind() Kind
		Mod() *Mod
	}
	ModLookup[T mod] interface {
		All() []T
		Get(m T) (found T, ok bool)
		GetByID(modID ModID) (found T, ok bool)
		Has(m T) bool
		Len() int
		Remove(m T)
		Set(m T)
	}
	modLookup[T mod] struct {
		Lookup map[lookupID]T `json:"Lookup"`
	}
)

func NewModLookup[T mod]() ModLookup[T] {
	return &modLookup[T]{
		Lookup: make(map[lookupID]T),
	}
}

func (l *modLookup[T]) All() []T {
	s := make([]T, 0, len(l.Lookup))
	for _, m := range l.Lookup {
		s = append(s, m)
	}
	return s
}

func (l *modLookup[T]) Set(m T) {
	l.Lookup[l.newLookupID(m)] = m
}

func (l *modLookup[T]) Has(m T) bool {
	_, ok := l.Lookup[l.newLookupID(m)]
	return ok
}

func (l *modLookup[T]) Get(m T) (found T, ok bool) {
	found, ok = l.Lookup[l.newLookupID(m)]
	return
}

func (l *modLookup[T]) GetByID(modID ModID) (found T, ok bool) {
	for _, m := range l.Lookup {
		if m.ID() == modID {
			found = m
			ok = true
			break
		}
	}
	return
}

func (l *modLookup[T]) Remove(m T) {
	delete(l.Lookup, l.newLookupID(m))
}

func (l *modLookup[T]) Len() int {
	return len(l.Lookup)
}

func (l *modLookup[T]) newLookupID(m T) lookupID {
	return lookupID(fmt.Sprintf("%s.%s", m.Kind(), m.ID()))
}
