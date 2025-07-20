package utils

import "sync"

type SafeSlice[T any] struct {
	mu    sync.Mutex
	items []T
}

func NewSafeSlice[T any]() SafeSlice[T] {
	return SafeSlice[T]{
		items: make([]T, 0),
	}
}

func (s *SafeSlice[T]) Append(val T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, val)
}

func (s *SafeSlice[T]) Items() []T {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]T, len(s.items))
	copy(out, s.items)
	return out
}

func (s *SafeSlice[T]) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.items)
}
