package cache

import (
    "errors"
    "sync"
    "testing"
    "time"
)

// TestCache_SetWithExpiration_Expired tests setting an item with expiration that correctly expires.
func TestCache_SetWithExpiration_Expired(t *testing.T) {
    cache := NewCache(10*time.Millisecond, 0)
    defer cache.StopJanitor()

    key := "expiringKey"
    value := "expiringValue"
    expiration := 5 * time.Millisecond
    cache.Set(key, value, expiration)

    time.Sleep(6 * time.Millisecond) // Ensure the item has expired

    _, err := cache.Get(key)
    if !errors.Is(err, ErrItemExpired) {
        t.Errorf("Expected error %v for expired item, got %v", ErrItemExpired, err)
    }
}

// TestCache_SetWithExpiration_NotExpired tests setting an item with expiration that has not yet expired.
func TestCache_SetWithExpiration_NotExpired(t *testing.T) {
    cache := NewCache(10*time.Millisecond, 0)
    defer cache.StopJanitor()

    key := "willExpireKey"
    value := "willExpireValue"
    expiration := 50 * time.Millisecond
    cache.Set(key, value, expiration)

    time.Sleep(1 * time.Millisecond) // Ensure the item has not expired yet

    got, err := cache.Get(key)
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if got != value {
        t.Errorf("Expected value %v, got %v", value, got)
    }
}

// TestCache_SetAndGetConcurrently_MultipleGoroutines tests concurrent setting and getting of items with multiple goroutines.
func TestCache_SetAndGetConcurrently_MultipleGoroutines(t *testing.T) {
    cache := NewCache(100*time.Millisecond, 0)
    defer cache.StopJanitor()

    var wg sync.WaitGroup
    key := "concurrentKey"
    value := "concurrentValue"

    // Set value in one goroutine
    wg.Add(1)
    go func() {
        defer wg.Done()
        cache.Set(key, value, 0)
    }()

    // Get value in another goroutine
    wg.Add(1)
    go func() {
        defer wg.Done()
        time.Sleep(1 * time.Millisecond) // Give some time for the Set operation
        got, err := cache.Get(key)
        if err != nil {
            t.Fatalf("Concurrent Get() unexpected error: %v", err)
        }
        if got != value {
            t.Errorf("Concurrent Get() = %v, want %v", got, value)
        }
    }()

    wg.Wait()
}

// TestCache_DeleteExisting_Success tests deleting an existing item.
func TestCache_DeleteExisting_Success(t *testing.T) {
    cache := NewCache(100*time.Millisecond, 0)
    defer cache.StopJanitor()

    key := "existentKey"
    value := "existentValue"
    cache.Set(key, value, 0)

    cache.Delete(key)

    _, err := cache.Get(key)
    if !errors.Is(err, ErrItemNotFound) {
        t.Errorf("Expected %v error, got %v", ErrItemNotFound, err)
    }
}

// TestCache_Clear_EmptyCache tests clearing an already empty cache.
func TestCache_Clear_EmptyCache(t *testing.T) {
    cache := NewCache(100*time.Millisecond, 0)
    defer cache.StopJanitor()

    cache.Clear() // Clear an empty cache

    if len(cache.items) != 0 {
        t.Errorf("Clear() on empty cache failed, expected cache to be empty, got %d items", len(cache.items))
    }
}

// TestCache_SetWithNegativeExpiration tests setting an item with negative expiration.
func TestCache_SetWithNegativeExpiration(t *testing.T) {
    cache := NewCache(100*time.Millisecond, 0)
    defer cache.StopJanitor()

    key := "negativeExpirationKey"
    value := "value"
    expiration := -5 * time.Millisecond
    cache.Set(key, value, expiration)

    _, err := cache.Get(key)
    if !errors.Is(err, ErrItemNotFound) {
        t.Errorf("Expected %v error for item with negative expiration, got %v", ErrItemNotFound, err)
    }
}