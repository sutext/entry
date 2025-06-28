package safe

import "sync"

type Any[T any] struct {
	mu sync.RWMutex
	v  T
}

func New[T any](v T) *Any[T] {
	return &Any[T]{v: v}
}
func (l *Any[T]) Get() T {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.v
}

func (l *Any[T]) Set(v T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.v = v
}
func (l *Any[T]) Write(fn func(T) T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.v = fn(l.v)
}
func (l *Any[T]) Read(fn func(T)) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	fn(l.v)
}
