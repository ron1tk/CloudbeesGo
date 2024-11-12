package cache

import (
	"testing"
	"time"
)

func setupCache(defaultDuration, cleanupInterval time.Duration, maxEntries int) *Cache {
	return NewCache(cleanupInterval, defaultDuration, maxEntries)
}

func TestNewCache(t *testing.T) {
	cache := setupCache(5*time.Minute, 1*time.Minute, 100)
	if cache == nil {
		t.Errorf("NewCache() returned nil")
	}
}

func TestCache_SetAndGet(t *testing.T) {
	cache := setupCache(5*time.Minute, 1*time.Minute, 10)

	cache.Set("key1", "value1", 0)
	val, err := cache.Get("key1")

	if err != nil || val != "value1" {
		t.Errorf("Set or Get failed. Err: %v, Val: %v", err, val)
	}

	// Test expiration
	cache.Set("keyExpire", "valueExpire", 1*time.Nanosecond)
	time.Sleep(2 * time.Nanosecond)
	_, err = cache.Get("keyExpire")

	if err != ErrItemExpired {
		t.Errorf("Expected ErrItemExpired, got: %v", err)
	}

	// Test non-existing key
	_, err = cache.Get("nonExistingKey")
	if err != ErrItemNotFound {
		t.Errorf("Expected ErrItemNotFound, got: %v", err)
	}
}

func TestCache_Update(t *testing.T) {
	cache := setupCache(5*time.Minute, 1*time.Minute, 10)
	cache.Set("key1", "value1", 0)

	err := cache.Update("key1", "newValue", 0)
	if err != nil {
		t.Errorf("Update failed. Err: %v", err)
	}

	val, _ := cache.Get("key1")
	if val != "newValue" {
		t.Errorf("Update did not change the value. Expected newValue, got: %v", val)
	}

	err = cache.Update("nonExistingKey", "value", 0)
	if err != ErrItemNotFound {
		t.Errorf("Expected ErrItemNotFound for non-existing key update, got: %v", err)
	}
}

func TestCache_Delete(t *testing.T) {
	cache := setupCache(5*time.Minute, 1*time.Minute, 10)
	cache.Set("key1", "value1", 0)

	cache.Delete("key1")
	_, err := cache.Get("key1")

	if err != ErrItemNotFound {
		t.Errorf("Delete did not remove the item. Err: %v", err)
	}
}

func TestCache_Exists(t *testing.T) {
	cache := setupCache(5*time.Minute, 1*time.Minute, 10)
	cache.Set("key1", "value1", 0)

	if !cache.Exists("key1") {
		t.Errorf("Exists reported false for existing key")
	}

	if cache.Exists("nonExistingKey") {
		t.Errorf("Exists reported true for non-existing key")
	}
}

func TestCache_Clear(t *testing.T) {
	cache := setupCache(5*time.Minute, 1*time.Minute, 10)
	cache.Set("key1", "value1", 0)
	cache.Set("key2", "value2", 0)

	cache.Clear()

	if len(cache.Keys()) != 0 {
		t.Errorf("Clear did not remove all items")
	}
}

func TestCache_Keys(t *testing.T) {
	cache := setupCache(5*time.Minute, 1*time.Minute, 10)
	cache.Set("key1", "value1", 0)
	cache.Set("key2", "value2", 0)

	keys := cache.Keys()
	if len(keys) != 2 {
		t.Errorf("Keys did not return correct number of keys. Expected 2, got: %d", len(keys))
	}
}

func TestCache_Stats(t *testing.T) {
	cache := setupCache(5*time.Minute, 1*time.Minute, 10)
	cache.Set("key1", "value1", 0)

	stats := cache.Stats()
	if stats.Items != 1 {
		t.Errorf("Stats did not return correct items count. Expected 1, got: %d", stats.Items)
	}
}

func TestCache_Eviction(t *testing.T) {
	cache := setupCache(5*time.Minute, 1*time.Minute, 1) // maxEntries set to 1
	cache.Set("key1", "value1", 0)
	cache.Set("key2", "value2", 0) // This should cause key1 to be evicted

	if cache.Exists("key1") {
		t.Errorf("LRU eviction failed. 'key1' should have been evicted.")
	}
}