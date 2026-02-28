package cache

import (
	"fmt"
	"time"
)

type Cache[T any] interface {
	Get(key string) (T, error)
	Set(key string, value T, ttl time.Duration) error
	Delete(key string) error
}

type memoryCache[T any] struct {
	cache map[string]T
}

func NewMemory[T any]() Cache[T] {
	return &memoryCache[T]{
		cache: make(map[string]T),
	}
}

func (c *memoryCache[T]) Get(key string) (T, error) {
	value, ok := c.cache[key]
	if !ok {
		var zero T
		return zero, fmt.Errorf("key %s not found", key)
	}
	return value, nil
}

func (c *memoryCache[T]) Set(key string, value T, ttl time.Duration) error {
	c.cache[key] = value
	return nil
}
func (c *memoryCache[T]) Delete(key string) error {
	delete(c.cache, key)
	return nil
}
