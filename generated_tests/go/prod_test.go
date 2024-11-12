package cache

import (
    "reflect"
    "testing"
    "time"
)

func TestNewCache(t *testing.T) {
    cleanupInterval := 1 * time.Millisecond
    cache := NewCache(cleanupInterval)
    defer cache.StopJanitor()

    if cache == nil {
        t.Error("NewCache() should not return nil")
    }

    if len(cache.items) != 0 {
        t.Errorf("New cache should be empty, got %d items", len(cache.items))
    }
}

func TestCache_SetAndGet(t *testing.T) {
    cache := NewCache(100 * time.Millisecond)
    defer cache.StopJanitor()

    key := "key"
    value := "value"
    cache.Set(key, value, 0) // No expiration

    got, err := cache.Get(key)
    if err != nil {
        t.Errorf("Get() unexpected error: %v", err)
    }
    if !reflect.DeepEqual(got, value) {
        t.Errorf("Get() = %v, want %v", got, value)
    }
}

func TestCache_Get_ItemNotFound(t *testing.T) {
    cache := NewCache(100 * time.Millisecond)
    defer cache.StopJanitor()

    _, err := cache.Get("nonexistent")
    if err == nil {
        t.Error("Get() expected error for nonexistent item")
    }
}

func TestCache_Get_ItemExpired(t *testing.T) {
    cache := NewCache(1 * time.Millisecond)
    defer cache.StopJanitor()

    key := "key"
    cache.Set(key, "value", 1*time.Nanosecond)

    time.Sleep(2 * time.Millisecond) // Ensure expiration

    _, err := cache.Get(key)
    if err == nil {
        t.Error("Get() expected error for expired item")
    }
}

func TestCache_Delete(t *testing.T) {
    cache := NewCache(100 * time.Millisecond)
    defer cache.StopJanitor()

    key := "key"
    cache.Set(key, "value", 0)
    cache.Delete(key)

    _, err := cache.Get(key)
    if err == nil {
        t.Errorf("Delete() failed, item %s still exists", key)
    }
}

func TestCache_DeleteExpired(t *testing.T) {
    cache := NewCache(1 * time.Millisecond)
    defer cache.StopJanitor()

    key := "key"
    cache.Set(key, "value", 1*time.Nanosecond)

    time.Sleep(2 * time.Millisecond) // Wait for item to expire

    cache.DeleteExpired()

    if _, exists := cache.items[key]; exists {
        t.Errorf("DeleteExpired() failed, expired item %s still exists", key)
    }
}

func TestCache_StopJanitor(t *testing.T) {
    cache := NewCache(1 * time.Millisecond)

    // Not an ideal test, just ensures calling StopJanitor doesn't result in panic
    cache.StopJanitor()
}