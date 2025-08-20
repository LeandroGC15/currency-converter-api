package cache

import (
	"sync"
	"time"
)

type item struct {
	value      interface{}
	expiration int64
}

type TTLCache struct {
	cache map[string]item
	mu    sync.RWMutex
}

func NewTTLCache() *TTLCache {
	return &TTLCache{
		cache: make(map[string]item),
	}
}

func (c *TTLCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = item{
		value:      value,
		expiration: time.Now().Add(ttl).UnixNano(),
	}
}

func (c *TTLCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.cache[key]
	if !found {
		return nil, false
	}

	if time.Now().UnixNano() > item.expiration {
		return nil, false
	}

	return item.value, true
}

func (c *TTLCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, key)
}
