package cache

import (
	"testing"
	"time"
)

func TestCache_SetWithNegativeExpiration(t *testing.T) {
	cache := setupCache(5*time.Minute, 1*time.Minute, 10)

	err := cache.Set("keyNegExp", "value", -1*time.Second)
	if err == nil {
		t.Errorf("Set with negative expiration did not fail")
	}
}

func TestCache_GetAfterCleanup(t *testing.T) {
	cache := setupCache(1*time.Nanosecond, 1*time.Nanosecond, 10)
	cache.Set("keyCleanup", "valueCleanup", 0)

	time.Sleep(2 * time.Nanosecond) // ensure cleanup has run
	_, err := cache.Get("keyCleanup")
	if err != ErrItemNotFound {
		t.Errorf("Expected ErrItemNotFound after cleanup, got: %v", err)
	}
}

func TestCache_UpdateExpiration(t *testing.T) {
	cache := setupCache(5*time.Minute, 1*time.Minute, 10)
	cache.Set("keyUpdateExp", "value", 1*time.Nanosecond)

	err := cache.Update("keyUpdateExp", "valueUpdated", 1*time.Hour)
	if err != nil {
		t.Errorf("Update expiration failed. Err: %v", err)
	}

	time.Sleep(2 * time.Nanosecond) // Past initial expiration
	_, err = cache.Get("keyUpdateExp")
	if err != nil {
		t.Errorf("Item should not have expired after update. Err: %v", err)
	}
}

func TestCache_DeleteNonExistingKey(t *testing.T) {
	cache := setupCache(5*time.Minute, 1*time.Minute, 10)

	err := cache.Delete("nonExistingKey")
	if err != ErrItemNotFound {
		t.Errorf("Expected ErrItemNotFound on deleting non-existing key, got: %v", err)
	}
}

func TestCache_ClearOnEmptyCache(t *testing.T) {
	cache := setupCache(5*time.Minute, 1*time.Minute, 10)

	cache.Clear() // Clearing an already empty cache

	if len(cache.Keys()) != 0 {
		t.Errorf("Clear on an empty cache did not behave as expected")
	}
}

func TestCache_StatsAfterDelete(t *testing.T) {
	cache := setupCache(5*time.Minute, 1*time.Minute, 10)
	cache.Set("key1", "value1", 0)
	cache.Delete("key1")

	stats := cache.Stats()
	if stats.Items != 0 {
		t.Errorf("Stats did not return correct items count after delete. Expected 0, got: %d", stats.Items)
	}
}

func TestCache_SetMaxEntriesZero(t *testing.T) {
	cache := setupCache(5*time.Minute, 1*time.Minute, 0) // maxEntries set to 0, should not limit entries

	cache.Set("key1", "value1", 0)
	cache.Set("key2", "value2", 0)

	if len(cache.Keys()) != 2 {
		t.Errorf("Setting maxEntries to 0 did not behave as unlimited. Keys count: %d", len(cache.Keys()))
	}
}

func TestCache_EvictionOrder(t *testing.T) {
	cache := setupCache(5*time.Minute, 1*time.Minute, 2) // maxEntries set to 2
	cache.Set("key1", "value1", 0)
	time.Sleep(1 * time.Nanosecond) // ensure different timestamps
	cache.Set("key2", "value2", 0)
	cache.Set("key3", "value3", 0) // This should cause key1 to be evicted

	if cache.Exists("key1") {
		t.Errorf("Eviction order incorrect. 'key1' should have been evicted first.")
	}
	if !cache.Exists("key2") || !cache.Exists("key3") {
		t.Errorf("Eviction order incorrect. 'key2' and 'key3' should exist.")
	}
}

func TestCache_ConcurrentAccess(t *testing.T) {
	cache := setupCache(5*time.Minute, 1*time.Minute, 100)

	// Simulate concurrent access
	go func() {
		for i := 0; i < 100; i++ {
			cache.Set(time.Now().String(), "value", 0)
		}
	}()
	go func() {
		for i := 0; i < 100; i++ {
			cache.Set(time.Now().String(), "value", 0)
		}
	}()

	time.Sleep(1 * time.Second) // Wait for goroutines to finish

	if len(cache.Keys()) != 200 {
		t.Errorf("Concurrent access did not result in correct number of keys. Expected 200, got: %d", len(cache.Keys()))
	}
}

func TestCache_MultipleDeletes(t *testing.T) {
	cache := setupCache(5*time.Minute, 1*time.Minute, 10)
	cache.Set("key1", "value1", 0)
	cache.Set("key2", "value2", 0)

	cache.Delete("key1")
	cache.Delete("key2")

	if cache.Exists("key1") || cache.Exists("key2") {
		t.Errorf("Multiple deletes did not remove the items as expected")
	}
}