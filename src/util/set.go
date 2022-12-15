package util

import "sync"

//algebraic set in go
type Set[T comparable] struct {
	mx *sync.RWMutex
	m  map[T]struct{}
}

//creates empty set of type T
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{m: make(map[T]struct{}), mx: &sync.RWMutex{}}
}

//creates set with given data (includes only unique values)
func NewSetWithData[T comparable](data ...T) *Set[T] {
	s := &Set[T]{m: make(map[T]struct{}), mx: &sync.RWMutex{}}
	s.Insert(data...)
	return s
}

//insert unique values to set
func (s *Set[T]) Insert(i ...T) {
	s.mx.Lock()
	for _, v := range i {
		s.m[v] = struct{}{}
	}
	s.mx.Unlock()
}

//remove the given values from set
func (s *Set[T]) Erase(i ...T) {
	s.mx.Lock()
	for _, v := range i {
		delete(s.m, v)
	}
	s.mx.Unlock()
}

//returns true if set has element i
func (s *Set[T]) Has(i T) bool {
	s.mx.RLock()
	defer s.mx.RUnlock()
	_, ok := s.m[i]
	return ok
}

//returns true if set has all given values
func (s *Set[T]) HasAll(i ...T) bool {
	s.mx.RLock()
	defer s.mx.RUnlock()
	var ok = true
	for v := 0; v < len(i) && ok; v++ {
		_, ok = s.m[i[v]]
	}
	return ok
}

//returns true if set has any of given values
func (s *Set[T]) HasAny(i ...T) bool {
	s.mx.RLock()
	defer s.mx.RUnlock()
	var ok = false
	for v := 0; v < len(i) && !ok; v++ {
		_, ok = s.m[i[v]]
	}
	return ok
}

//returns the number of elements in set
func (s *Set[T]) Size() int {
	s.mx.RLock()
	defer s.mx.RUnlock()
	var cnt int
	for range s.m {
		cnt++
	}
	return cnt
}

//returns the list of set's elements
func (s *Set[T]) ToList() []T {
	s.mx.RLock()
	defer s.mx.RUnlock()
	l := make([]T, 0)
	for v := range s.m {
		l = append(l, v)
	}
	return l
}

//returns the intersection of sets
func (s *Set[T]) Intersect(b *Set[T]) *Set[T] {
	s.mx.RLock()
	b.mx.RLock()
	defer s.mx.RUnlock()
	defer b.mx.RUnlock()
	c := NewSet[T]()
	for k1 := range s.m {
		if b.Has(k1) {
			c.Insert(k1)
		}
	}
	return c
}
