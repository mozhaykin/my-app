package dto

import (
	"sync"
)

type Cache struct {
	mx sync.RWMutex
	m  map[string]int
}

func NewCache() *Cache {
	return &Cache{
		m: make(map[string]int, 100),
	}
}

func (c *Cache) Add(name string, value int) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.m[name] = value
}

func (c *Cache) Get(name string) (int, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	value, ok := c.m[name]

	return value, ok
}
