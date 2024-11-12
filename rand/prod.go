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
    items   map[string]Item
    mutex   sync.RWMutex
    janitor *janitor
}

// NewCache creates a new Cache instance and starts the janitor.
func NewCache(cleanupInterval time.Duration) *Cache {
    c := &Cache{
        items: make(map[string]Item),
    }
    runJanitor(c, cleanupInterval)
    return c
}

// Set adds an item to the cache with a specified duration.
// If duration is 0, the item does not expire.
func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
    var expiration int64
    if duration > 0 {
        expiration = time.Now().Add(duration).UnixNano()
    } else {
        expiration = 0
    }

    c.mutex.Lock()
    defer c.mutex.Unlock()
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
        return nil, errors.New("item not found")
    }

    if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
        c.mutex.Lock()
        delete(c.items, key)
        c.mutex.Unlock()
        return nil, errors.New("item expired")
    }

    return item.Value, nil
}

// Delete removes an item from the cache.
func (c *Cache) Delete(key string) {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    delete(c.items, key)
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

// DeleteExpired removes all expired items from the cache.
func (c *Cache) DeleteExpired() {
    now := time.Now().UnixNano()
    c.mutex.Lock()
    defer c.mutex.Unlock()
    for k, v := range c.items {
        if v.Expiration > 0 && now > v.Expiration {
            delete(c.items, k)
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

// StopJanitor stops the janitor goroutine.
func (c *Cache) StopJanitor() {
    c.janitor.stop <- true
}