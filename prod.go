package cache

import (
	"errors"
	"sync"
	"time"
)

// Item represents a single cache item.
type Item struct {
	Value      interface{}
	Expiration int64
}

// Cache represents the in-memory cache.
type Cache struct {
	items           map[string]Item
	mutex           sync.RWMutex
	janitor         *janitor
	defaultDuration time.Duration
	stats           CacheStats
	evictionCallback func(key string, value interface{})
}

// CacheStats holds statistics about cache usage.
type CacheStats struct {
	Hits   int
	Misses int
	Items  int
}

// NewCache creates a new Cache instance and starts the janitor.
// If defaultDuration is 0, items will not expire unless a specific duration is set.
func NewCache(cleanupInterval time.Duration, defaultDuration time.Duration) *Cache {
	c := &Cache{
		items:           make(map[string]Item),
		defaultDuration: defaultDuration,
	}
	runJanitor(c, cleanupInterval)
	return c
}

// Set adds an item to the cache with a specified duration.
// If duration is 0, the default duration is used.
// If both are 0, the item does not expire.
func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	var expiration int64
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	} else if c.defaultDuration > 0 {
		expiration = time.Now().Add(c.defaultDuration).UnixNano()
	} else {
		expiration = 0
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, exists := c.items[key]
	if !exists {
		c.stats.Items++
	}
	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
	}
}

// Get retrieves an item from the cache.
// Returns an error if the item does not exist or has expired.
func (c *Cache) Get(key string) (interface{}, error) {
	c.mutex.RLock()
	item, found := c.items[key]
	c.mutex.RUnlock()

	if !found {
		c.incrementMisses()
		return nil, errors.New("item not found")
	}

	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		c.mutex.Lock()
		delete(c.items, key)
		c.stats.Items--
		c.mutex.Unlock()
		c.incrementMisses()
		if c.evictionCallback != nil {
			c.evictionCallback(key, item.Value)
		}
		return nil, errors.New("item expired")
	}

	c.incrementHits()
	return item.Value, nil
}

// Delete removes an item from the cache.
func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	item, exists := c.items[key]
	if exists {
		delete(c.items, key)
		c.stats.Items--
	}
	c.mutex.Unlock()
	if exists && c.evictionCallback != nil {
		c.evictionCallback(key, item.Value)
	}
}

// Clear removes all items from the cache.
func (c *Cache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for k, v := range c.items {
		if c.evictionCallback != nil {
			c.evictionCallback(k, v.Value)
		}
	}
	c.items = make(map[string]Item)
	c.stats.Items = 0
}

// Exists checks if a key exists in the cache without retrieving its value.
func (c *Cache) Exists(key string) bool {
	c.mutex.RLock()
	item, found := c.items[key]
	c.mutex.RUnlock()

	if !found {
		return false
	}

	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		c.mutex.Lock()
		delete(c.items, key)
		c.stats.Items--
		c.mutex.Unlock()
		if c.evictionCallback != nil {
			c.evictionCallback(key, item.Value)
		}
		return false
	}

	return true
}

// Keys returns a slice of all keys currently stored in the cache.
func (c *Cache) Keys() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	keys := make([]string, 0, len(c.items))
	for k, v := range c.items {
		if v.Expiration == 0 || time.Now().UnixNano() <= v.Expiration {
			keys = append(keys, k)
		}
	}
	return keys
}

// Stats returns the current cache statistics.
func (c *Cache) Stats() CacheStats {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.stats
}

// SetEvictionCallback sets a callback function that is called whenever an item is evicted.
func (c *Cache) SetEvictionCallback(callback func(key string, value interface{})) {
	c.evictionCallback = callback
}

// StopJanitor stops the janitor goroutine.
func (c *Cache) StopJanitor() {
	c.janitor.stop <- true
}

// DeleteExpired removes all expired items from the cache.
func (c *Cache) DeleteExpired() {
	now := time.Now().UnixNano()
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for k, v := range c.items {
		if v.Expiration > 0 && now > v.Expiration {
			delete(c.items, k)
			c.stats.Items--
			if c.evictionCallback != nil {
				c.evictionCallback(k, v.Value)
			}
		}
	}
}

// incrementHits safely increments the hit counter.
func (c *Cache) incrementHits() {
	c.mutex.Lock()
	c.stats.Hits++
	c.mutex.Unlock()
}

// incrementMisses safely increments the miss counter.
func (c *Cache) incrementMisses() {
	c.mutex.Lock()
	c.stats.Misses++
	c.mutex.Unlock()
}

// janitor is responsible for cleaning up expired items.
type janitor struct {
	Interval time.Duration
	stop     chan bool
}

// Run starts the janitor to periodically clean up expired items.
func (j *janitor) Run(c *Cache) {
	ticker := time.NewTicker(j.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-j.stop:
			return
		}
	}
}

// runJanitor initializes and starts the janitor.
func runJanitor(c *Cache, ci time.Duration) {
	j := &janitor{
		Interval: ci,
		stop:     make(chan bool),
	}
	c.janitor = j
	go j.Run(c)
}
