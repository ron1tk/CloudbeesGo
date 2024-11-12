package cache

import (
	"reflect"
	"testing"
	"time"
)

func TestSetupCache(t *testing.T) {
	c := setupCache(1 * time.Minute)
	if c == nil {
		t.Fatal("setupCache() failed, got nil")
	}
}

func TestTeardownCache(t *testing.T) {
	c := setupCache(1 * time.Minute)
	if c == nil {
		t.Fatal("setupCache() returned nil, cannot test teardownCache()")
	}
	defer teardownCache(c)
	// Assuming teardownCache() modifies the cache in a testable way; otherwise, this is a placeholder
	if c.janitorTicker != nil {
		t.Error("teardownCache() should stop the janitor ticker")
	}
}

func TestCache_Set_ErrorCase(t *testing.T) {
	c := setupCache(1 * time.Minute)
	defer teardownCache(c)

	err := c.Set("", "value", 1*time.Hour)
	if err == nil {
		t.Error("Cache.Set() with empty key expected to return an error, got nil")
	}
}

func TestCache_Get_ErrorCase(t *testing.T) {
	c := setupCache(1 * time.Minute)
	defer teardownCache(c)

	_, err := c.Get("")
	if err == nil {
		t.Error("Cache.Get() with empty key expected to return an error, got nil")
	}
}

func TestCache_Delete_NonExistingKey(t *testing.T) {
	c := setupCache(1 * time.Minute)
	defer teardownCache(c)

	c.Delete("nonExistingKey")
	_, err := c.Get("nonExistingKey")
	if err == nil {
		t.Error("Cache.Delete() non-existing key expected to not find the item, but it did")
	}
}

func TestCache_TTLBehaviour(t *testing.T) {
	c := setupCache(1 * time.Minute)
	defer teardownCache(c)

	c.Set("keyWithShortTTL", "value", 1*time.Millisecond)
	time.Sleep(2 * time.Millisecond)

	_, err := c.Get("keyWithShortTTL")
	if err == nil {
		t.Error("Expected error for key with expired TTL, got nil")
	}
}

func TestCache_Concurrency(t *testing.T) {
	c := setupCache(1 * time.Minute)
	defer teardownCache(c)

	go func() {
		c.Set("concurrentKey1", "value1", 1*time.Hour)
	}()
	go func() {
		c.Set("concurrentKey2", "value2", 1*time.Hour)
	}()

	time.Sleep(1 * time.Second) // Wait for goroutines to complete

	v1, err1 := c.Get("concurrentKey1")
	if err1 != nil || v1 == nil {
		t.Error("Failed to get value for concurrentKey1")
	}

	v2, err2 := c.Get("concurrentKey2")
	if err2 != nil || v2 == nil {
		t.Error("Failed to get value for concurrentKey2")
	}
}

func TestCache_Clear(t *testing.T) {
	c := setupCache(1 * time.Minute)
	defer teardownCache(c)

	c.Set("keyToClear", "value", 1*time.Hour)
	c.Clear()

	_, err := c.Get("keyToClear")
	if err == nil {
		t.Errorf("Cache.Clear() failed, expected error for cleared key, got nil")
	}
}

func TestCache_MultipleOperations(t *testing.T) {
	c := setupCache(1 * time.Minute)
	defer teardownCache(c)

	// Set and immediately delete
	c.Set("tempKey", "tempValue", 1*time.Hour)
	c.Delete("tempKey")

	_, err := c.Get("tempKey")
	if err == nil {
		t.Errorf("Expected error after set and delete operation, got nil")
	}

	// Set and check existence
	c.Set("permanentKey", "value", 1*time.Hour)
	value, err := c.Get("permanentKey")
	if err != nil || value == nil {
		t.Errorf("Failed to get value for permanentKey")
	}
}