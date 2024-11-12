package cache

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

// Define exported errors for better error handling and testing
var (
	ErrItemNotFound = errors.New("item not found")
	ErrItemExpired  = errors.New("item expired")
)

// Item represents a single cache item.
type Item struct {
	Value      interface{}
	Expiration int64
}

// Cache represents the in-memory cache.
type Cache struct {
	items            map[string]Item
	mutex            sync.RWMutex
	janitor          *janitor
	defaultDuration  time.Duration
	stats            CacheStats
	evictionCallback func(key string, value interface{})

	// LRU-related fields
	maxEntries int
	lruList    *list.List                   // List to maintain LRU order
	lruMap     map[string]*list.Element     // Map to quickly access list elements
}

// CacheStats holds statistics about cache usage.
type CacheStats struct {
	Hits      int
	Misses    int
	Items     int
	Evictions int
}

// NewCache creates a new Cache instance and starts the janitor.
// If defaultDuration is 0, items will not expire unless a specific duration is set.
// If maxEntries is greater than 0, the cache will enforce a maximum number of items using LRU eviction.
func NewCache(cleanupInterval time.Duration, defaultDuration time.Duration, maxEntries int) *Cache {
	c := &Cache{
		items:           make(map[string]Item),
		defaultDuration: defaultDuration,
		maxEntries:      maxEntries,
		lruList:         list.New(),
		lruMap:          make(map[string]*list.Element),
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

	// If the item already exists, update it and move it to the front of the LRU list
	if element, exists := c.lruMap[key]; exists {
		c.lruList.MoveToFront(element)
		element.Value = key
	} else {
		// If adding a new item, check for capacity
		if c.maxEntries > 0 && c.lruList.Len() >= c.maxEntries {
			// Evict the least recently used item
			c.evictOldest()
		}
		// Add the new item to the front of the LRU list
		element := c.lruList.PushFront(key)
		c.lruMap[key] = element
		c.stats.Items++
	}

	// Set or update the item
	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
	}
}

// Get retrieves an item from the cache.
// Returns an error if the item does not exist or has expired.
func (c *Cache) Get(key string) (interface{}, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, found := c.items[key]
	if !found {
		c.incrementMisses()
		return nil, ErrItemNotFound
	}

	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		// Item has expired
		c.deleteItem(key)
		c.incrementMisses()
		return nil, ErrItemExpired
	}

	// Move the accessed item to the front of the LRU list
	if element, exists := c.lruMap[key]; exists {
		c.lruList.MoveToFront(element)
	}

	c.incrementHits()
	return item.Value, nil
}

// Update modifies the value and/or expiration of an existing item.
// Returns an error if the item does not exist or has expired.
func (c *Cache) Update(key string, value interface{}, duration time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, found := c.items[key]
	if !found {
		c.incrementMisses()
		return ErrItemNotFound
	}

	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		// Item has expired
		c.deleteItem(key)
		c.incrementMisses()
		return ErrItemExpired
	}

	// Update value and expiration
	var expiration int64
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	} else if c.defaultDuration > 0 {
		expiration = time.Now().Add(c.defaultDuration).UnixNano()
	} else {
		expiration = 0
	}

	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
	}

	// Move the updated item to the front of the LRU list
	if element, exists := c.lruMap[key]; exists {
		c.lruList.MoveToFront(element)
	}

	return nil
}

// Delete removes an item from the cache.
func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.deleteItem(key)
}

// deleteItem is a helper function to remove an item without locking.
// Assumes the caller holds the lock.
func (c *Cache) deleteItem(key string) {
	item, exists := c.items[key]
	if !exists {
		return
	}

	delete(c.items, key)
	if element, exists := c.lruMap[key]; exists {
		c.lruList.Remove(element)
		delete(c.lruMap, key)
	}
	c.stats.Items--

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
	c.lruList.Init()
	c.lruMap = make(map[string]*list.Element)
	c.stats.Items = 0
	c.stats.Evictions += c.stats.Items
}

// Exists checks if a key exists in the cache without retrieving its value.
func (c *Cache) Exists(key string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, found := c.items[key]
	if !found {
		return false
	}

	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		// Item has expired
		c.deleteItem(key)
		return false
	}

	// Move the accessed item to the front of the LRU list
	if element, exists := c.lruMap[key]; exists {
		c.lruList.MoveToFront(element)
	}

	return true
}

// Keys returns a slice of all keys currently stored in the cache.
func (c *Cache) Keys() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	keys := make([]string, 0, len(c.items))
	now := time.Now().UnixNano()
	for k, v := range c.items {
		if v.Expiration == 0 || now <= v.Expiration {
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
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.evictionCallback = callback
}

// StopJanitor stops the janitor goroutine.
func (c *Cache) StopJanitor() {
	if c.janitor != nil {
		c.janitor.stop <- true
	}
}

// DeleteExpired removes all expired items from the cache.
func (c *Cache) DeleteExpired() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	now := time.Now().UnixNano()
	for k, v := range c.items {
		if v.Expiration > 0 && now > v.Expiration {
			c.deleteItem(k)
			c.stats.Items--
			c.stats.Evictions++
		}
	}
}

// incrementHits safely increments the hit counter.
func (c *Cache) incrementHits() {
	c.stats.Hits++
}

// incrementMisses safely increments the miss counter.
func (c *Cache) incrementMisses() {
	c.stats.Misses++
}

// evictOldest removes the least recently used item from the cache.
func (c *Cache) evictOldest() {
	element := c.lruList.Back()
	if element != nil {
		key := element.Value.(string)
		c.deleteItem(key)
		c.stats.Evictions++
	}
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
