package abc

import "sync"

type stack struct {
	items []float64
	mutex sync.Mutex
}

func newStack() *stack {
	return &stack{make([]float64, 0), sync.Mutex{}}
}

func (s *stack) push(item float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.items = append(s.items, item)
}

func (s *stack) pop() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}
