package abc

import "sync"

type stack struct {
	values []value
	mutex  sync.Mutex
}

func newStack() *stack {
	return &stack{make([]value, 0), sync.Mutex{}}
}

func (s *stack) push(v value) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.values = append(s.values, v)
}

func (s *stack) pop() value {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	item := s.values[len(s.values)-1]
	s.values = s.values[:len(s.values)-1]
	return item
}

func (s *stack) peek(distance int) value {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.values[len(s.values)-1-distance]
}
