package cache

import (
    "testing"
    "time"
)

func TestCache_NewCache_NonNil(t *testing.T) {
    cleanupInterval := 50 * time.Millisecond
    cache := NewCache(cleanupInterval)
    defer cache.StopJanitor()

    if cache == nil {
        t.Fatal("NewCache() returned nil, expected non-nil Cache instance")
    }
}

func TestCache_NewCache_EmptyOnInit(t *testing.T) {
    cache := NewCache(50 * time.Millisecond)
    defer cache.StopJanitor()

    if len(cache.items) != 0 {
        t.Fatalf("New cache expected to be empty, got %d items", len(cache.items))
    }
}

func TestCache_SetAndGet_Success(t *testing.T) {
    cache := NewCache(100 * time.Millisecond)
    defer cache.StopJanitor()

    key := "testKey"
    expectedValue := "testValue"
    cache.Set(key, expectedValue, 0) // No expiration

    actualValue, err := cache.Get(key)
    if err != nil {
        t.Fatalf("Get() returned unexpected error: %v", err)
    }
    if actualValue != expectedValue {
        t.Errorf("Get() = %v, want %v", actualValue, expectedValue)
    }
}

func TestCache_Get_KeyNotFound(t *testing.T) {
    cache := NewCache(100 * time.Millisecond)
    defer cache.StopJanitor()

    _, err := cache.Get("nonexistentKey")
    if err == nil {
        t.Error("Expected an error for a nonexistent key, got nil")
    }
}

func TestCache_Get_KeyExpired(t *testing.T) {
    cache := NewCache(1 * time.Millisecond)
    defer cache.StopJanitor()

    key := "expiringKey"
    cache.Set(key, "value", 1*time.Nanosecond)

    time.Sleep(2 * time.Millisecond) // Ensure the item is expired

    _, err := cache.Get(key)
    if err == nil {
        t.Error("Expected an error for an expired key, got nil")
    }
}

func TestCache_Delete_KeyExists(t *testing.T) {
    cache := NewCache(100 * time.Millisecond)
    defer cache.StopJanitor()

    key := "keyToDelete"
    cache.Set(key, "value", 0) // No expiration
    cache.Delete(key)

    _, err := cache.Get(key)
    if err == nil {
        t.Errorf("Expected an error after deleting key %s, got nil", key)
    }
}

func TestCache_DeleteExpired_ItemsExist(t *testing.T) {
    cache := NewCache(1 * time.Millisecond)
    defer cache.StopJanitor()

    key := "expiringKey"
    cache.Set(key, "value", 1*time.Nanosecond)

    time.Sleep(2 * time.Millisecond) // Ensure the item is expired

    cache.DeleteExpired()

    if _, exists := cache.items[key]; exists {
        t.Errorf("Expected expired item %s to be deleted, but it still exists", key)
    }
}

func TestCache_StopJanitor_NoPanic(t *testing.T) {
    cache := NewCache(1 * time.Millisecond)

    defer func() {
        if r := recover(); r != nil {
            t.Errorf("StopJanitor() caused panic: %v", r)
        }
    }()

    cache.StopJanitor()
}