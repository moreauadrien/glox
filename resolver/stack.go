// Package stack creates a ItemStack data structure for the Item type
package resolver

import (
	"sync"
)

// Stack the stack of Items
type Stack[T any] struct {
	items []T
	lock  sync.RWMutex
}

// New creates a new ItemStack
func (s *Stack[T]) New() *Stack[T] {
	s.items = []T{}
	return s
}

// Push adds an Item to the top of the stack
func (s *Stack[T]) Push(t T) {
	s.lock.Lock()
	s.items = append(s.items, t)
	s.lock.Unlock()
}

// Pop removes an Item from the top of the stack
func (s *Stack[T]) Pop() *T {
	s.lock.Lock()
	item := s.items[len(s.items)-1]
	s.items = s.items[0 : len(s.items)-1]
	s.lock.Unlock()
	return &item
}

func (s *Stack[T]) Peek() *T {
    return &s.items[len(s.items)-1]
}

func (s *Stack[T]) IsEmpty() bool {
    return len(s.items) == 0
}

func (s *Stack[T]) Len() int {
    return len(s.items)
}

func (s *Stack[T]) Get(i int) *T {
    return &s.items[i]
}
