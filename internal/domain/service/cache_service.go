package service

import (
	"sync"
	"time"
)

type Cache struct {
	DB         map[string]interface{}
	Expiration time.Duration
	mutex      sync.RWMutex
}

func NewCache(expiration time.Duration) *Cache {
	cache := &Cache{
		DB:         make(map[string]interface{}),
		Expiration: expiration,
	}
	go cache.cleanupLoop()
	return cache
}

func (c *Cache) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.DB[key] = value
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	value, exists := c.DB[key]
	return value, exists
}

func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.DB, key)
}

func (c *Cache) cleanupLoop() {
	ticker := time.NewTicker(c.Expiration)
	for range ticker.C {
		c.cleanup()
	}
}

func (c *Cache) cleanup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.DB = make(map[string]interface{})
}
