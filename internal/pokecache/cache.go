package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries      map[string]cacheEntry
	entriesMutex sync.Mutex
	interval     time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		entries:  make(map[string]cacheEntry),
		interval: interval,
	}

	cache.reapLoop()

	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.entriesMutex.Lock()
	defer c.entriesMutex.Unlock()

	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.entriesMutex.Lock()
	defer c.entriesMutex.Unlock()

	if entry, ok := c.entries[key]; ok {
		return entry.val, true
	}

	return nil, false
}

func (c *Cache) reapLoop() {
	go func() {
		ticker := time.NewTicker(c.interval)
		for range ticker.C {
			now := time.Now()
			c.entriesMutex.Lock()
			for key, entry := range c.entries {
				if now.Sub(entry.createdAt) > (c.interval - 1) {
					delete(c.entries, key)
				}
			}
			c.entriesMutex.Unlock()
		}
	}()
}
