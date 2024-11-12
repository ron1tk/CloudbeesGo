package cache

import (
    "reflect"
    "testing"
    "time"
)

func setupCache(cleanupInterval time.Duration) *Cache {
    return NewCache(cleanupInterval)
}

func teardownCache(c *Cache) {
    c.StopJanitor()
}

func TestCache_SetAndGet(t *testing.T) {
    c := setupCache(1 * time.Minute)
    defer teardownCache(c)

    c.Set("key1", "value1", 1*time.Hour)

    value, err := c.Get("key1")
    if err != nil {
        t.Errorf("Get() error = %v, wantErr %v", err, false)
    }
    if !reflect.DeepEqual(value, "value1") {
        t.Errorf("Get() = %v, want %v", value, "value1")
    }
}

func TestCache_Get_ItemNotFound(t *testing.T) {
    c := setupCache(1 * time.Minute)
    defer teardownCache(c)

    _, err := c.Get("nonexistent")
    if err == nil {
        t.Errorf("Get() error = %v, wantErr %v", err, true)
    }
}

func TestCache_Get_ItemExpired(t *testing.T) {
    c := setupCache(1 * time.Millisecond)
    defer teardownCache(c)

    c.Set("key2", "value2", 1*time.Nanosecond)

    time.Sleep(2 * time.Millisecond) // Wait for item to expire

    _, err := c.Get("key2")
    if err == nil {
        t.Errorf("Get() should error for expired item, but got %v", err)
    }
}

func TestCache_Delete(t *testing.T) {
    c := setupCache(1 * time.Minute)
    defer teardownCache(c)

    c.Set("key3", "value3", 1*time.Hour)
    c.Delete("key3")

    _, err := c.Get("key3")
    if err == nil {
        t.Errorf("Get() after Delete() should error, but got nil")
    }
}

func TestCache_DeleteExpired(t *testing.T) {
    c := setupCache(1 * time.Millisecond)
    defer teardownCache(c)

    c.Set("key4", "value4", 1*time.Nanosecond)
    time.Sleep(2 * time.Millisecond) // Wait for item to expire

    c.DeleteExpired()

    _, found := c.items["key4"]
    if found {
        t.Errorf("DeleteExpired() did not delete the expired item.")
    }
}

func TestNewCache(t *testing.T) {
    c := setupCache(1 * time.Minute)
    defer teardownCache(c)

    if c == nil {
        t.Errorf("NewCache() returned nil")
    }
}

func TestStopJanitor(t *testing.T) {
    c := setupCache(1 * time.Millisecond)

    // This test ensures calling StopJanitor does not produce panic or error
    // We do not have a direct way to assert the janitor has stopped,
    // but not crashing is a good sign of proper handling.
    c.StopJanitor()
}