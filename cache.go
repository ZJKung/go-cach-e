package gocache

import (
	"sync"

	"zjkung.github.com/g-cach-e/lru"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func NewCache(maxBytes int64) *cache {
	return &cache{
		mu:         sync.Mutex{},
		lru:        lru.NewLru(maxBytes, nil),
		cacheBytes: maxBytes,
	}
}
func (c *cache) Add(key string, value ReadOnlyByte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lru.Add(key, value)
}

func (c *cache) Get(key string) (value ReadOnlyByte, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ReadOnlyByte), ok
	}
	return
}
