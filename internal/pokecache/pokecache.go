package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries  map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
	stop     chan struct{}
}

func NewCache(wait time.Duration) *Cache {
	fmt.Println("creating cache")
	cache := &Cache{
		entries:  make(map[string]cacheEntry),
		mu:       sync.Mutex{},
		interval: wait,
		stop:     make(chan struct{}),
	}
	go cache.reapLoop(wait)
	return cache
}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now().UTC(),
		val:       value,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	cacheE, ok := c.entries[key]
	return cacheE.val, ok
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.reap(interval)
		case <-c.stop:
			// Received stop signal, exit the loop
			return
		}
	}
}

func (c *Cache) reap(interval time.Duration) {
	timeAgo := time.Now().UTC().Add(-interval)
	for k, v := range c.entries {
		if v.createdAt.Before(timeAgo) {
			delete(c.entries, k)
		}
	}
}

func (c *Cache) StopReapLoop() {
	close(c.stop)
}
