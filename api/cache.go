package api

import (
	"sync"
	"time"
)

// CacheEntry holds a cached search result with an expiration time.
type CacheEntry struct {
	Results   []map[string]interface{}
	ExpiresAt time.Time
}

// SearchCache is a simple in-memory cache for search results.
type SearchCache struct {
	mu      sync.RWMutex
	entries map[string]*CacheEntry
	ttl     time.Duration
}

// NewSearchCache creates a new SearchCache with the given TTL.
func NewSearchCache(ttl time.Duration) *SearchCache {
	return &SearchCache{
		entries: make(map[string]*CacheEntry),
		ttl:     ttl,
	}
}

// Get retrieves a cached result by key. Returns nil if not found or expired.
func (c *SearchCache) Get(key string) ([]map[string]interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}
	return entry.Results, true
}

// Set stores a result in the cache under the given key.
func (c *SearchCache) Set(key string, results []map[string]interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = &CacheEntry{
		Results:   results,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// Delete removes a cache entry by key.
func (c *SearchCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Flush removes all expired entries from the cache.
func (c *SearchCache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.entries {
		if now.After(entry.ExpiresAt) {
			delete(c.entries, key)
		}
	}
}

// Len returns the number of entries currently in the cache (including expired).
func (c *SearchCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
