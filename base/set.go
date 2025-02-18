package base

import (
	"sync"
)

// 去重集合
type Set[T comparable] struct {
	m map[T]struct{}

	lock sync.Mutex
}

func NewSetOfSize[T comparable](size int) *Set[T] {
	if size == 0 {
		size = 8
	}

	return &Set[T]{m: make(map[T]struct{}, size)}
}

func NewSet[T comparable](keys ...T) *Set[T] {
	s := NewSetOfSize[T](len(keys))

	for _, k := range keys {
		s.Add(k)
	}

	return s
}

func (s *Set[T]) Add(keys ...T) {
	s.init()

	for _, k := range keys {
		s.m[k] = struct{}{}
	}
}

func (s *Set[T]) Remove(key T) {
	s.init()

	delete(s.m, key)
}

func (s *Set[T]) Clear() {
	s.init()

	for k := range s.m {
		delete(s.m, k)
	}
}

func (s *Set[T]) Contains(key T) bool {
	s.init()

	if _, ok := s.m[key]; ok {
		return true
	} else {
		return false
	}
}

func (s *Set[T]) Get(keys ...T) (T, bool) {
	s.init()

	for _, k := range keys {
		if _, ok := s.m[k]; ok {
			return k, true
		}
	}

	for k := range s.m {
		return k, true
	}

	var k T
	return k, false
}

func (s *Set[T]) Range(f func(key T) bool) {
	s.init()

	for k := range s.m {
		if !f(k) {
			break
		}
	}
}

func (s *Set[T]) Range2(f func(key T) (bool, error)) error {
	s.init()

	for k := range s.m {
		ok, err := f(k)
		if err != nil {
			return err
		}

		if !ok {
			break
		}
	}

	return nil
}

func (s *Set[T]) Size() int {
	s.init()

	return len(s.m)
}

func (s *Set[T]) ToSlice() []T {
	s.init()

	ss := []T{}
	for k := range s.m {
		ss = append(ss, k)
	}
	return ss
}

func (s *Set[T]) init() {
	if s.m != nil {
		return
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	if s.m != nil {
		return
	}

	s.m = make(map[T]struct{})
}
