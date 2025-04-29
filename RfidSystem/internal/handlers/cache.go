package handlers

import (
	"container/list"
	"sync"
	"time"
)

// CacheItem represents a single cache entry
// V is the value type
// K is the key type
// For simplicity, we'll use string as the key type

type cacheItem struct {
	key       string
	value     any
	expiresAt time.Time
}

// LRUCache represents an LRU cache with a time-to-live (TTL) for cache items.
// It is safe for concurrent use since I embedded a Mutex in it.
type LRUCache struct {
	capacity int
	ttl      time.Duration
	mu       sync.Mutex
	items    map[string]*list.Element
	order    *list.List // Front is most recent
}

// NewLRUCache creates a new LRUCache with the specified capacity and time-to-live (TTL).
func NewLRUCache(capacity int, ttl time.Duration) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		ttl:      ttl,
		items:    make(map[string]*list.Element),
		order:    list.New(),
	}
}

// Get retrieves a value from the cache for a given key.
// It returns the value and a boolean indicating if the key was found and the item was not expired.
func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	elem, ok := c.items[key]
	if !ok {
		return nil, false
	}
	item := elem.Value.(*cacheItem)
	if time.Now().After(item.expiresAt) {
		// Expired, remove
		c.order.Remove(elem)
		delete(c.items, key)
		return nil, false
	}
	// Move to front
	c.order.MoveToFront(elem)
	return item.value, true
}

// Set adds or updates a key-value pair in the cache.
// If the cache exceeds its capacity, the least recently used item is removed.
func (c *LRUCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, ok := c.items[key]; ok {
		// Update existing
		item := elem.Value.(*cacheItem)
		item.value = value
		item.expiresAt = time.Now().Add(c.ttl)
		c.order.MoveToFront(elem)
		return
	}
	// Add new
	item := &cacheItem{
		key:       key,
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
	elem := c.order.PushFront(item)
	c.items[key] = elem
	if c.order.Len() > c.capacity {
		// Remove least recently used
		oldest := c.order.Back()
		if oldest != nil {
			c.order.Remove(oldest)
			oldestItem := oldest.Value.(*cacheItem)
			delete(c.items, oldestItem.key)
		}
	}
}

// Optional: Delete key
// Delete removes a key-value pair from the cache.
func (c *LRUCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, ok := c.items[key]; ok {
		c.order.Remove(elem)
		delete(c.items, key)
	}
}

// PurgeExpired removes all expired items from the cache.
func (c *LRUCache) PurgeExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for key, elem := range c.items {
		item := elem.Value.(*cacheItem)
		if now.After(item.expiresAt) {
			c.order.Remove(elem)
			delete(c.items, key)
		}
	}
}
