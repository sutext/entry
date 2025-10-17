package safe

import "sync"

type Set[Value comparable] struct {
	mu sync.RWMutex
	mp map[Value]struct{}
}

func NewSet[Value comparable](v ...Value) *Set[Value] {
	switch len(v) {
	case 0:
		return &Set[Value]{
			mp: make(map[Value]struct{}),
		}
	case 1:
		return &Set[Value]{
			mp: map[Value]struct{}{v[0]: {}},
		}
	default:
		panic("too many arguments")
	}
}
func (s *Set[Value]) Has(value Value) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.mp[value]
	return ok
}
func (s *Set[Value]) Add(value Value) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mp[value] = struct{}{}
}

func (s *Set[Value]) Remove(value Value) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.mp, value)
}
func (s *Set[Value]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.mp)
}
