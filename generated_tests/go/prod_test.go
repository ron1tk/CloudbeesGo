package cache

import (
    "testing"
    "time"
)

// TestCache_SetWithExpiration_Success tests setting an item with expiration successfully.
func TestCache_SetWithExpiration_Success(t *testing.T) {
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

// TestCache_SetAndGetConcurrently_Success tests concurrent setting and getting of items.
func TestCache_SetAndGetConcurrently_Success(t *testing.T) {
    cache := NewCache(100*time.Millisecond, 0)
    defer cache.StopJanitor()

    key := "concurrentKey"
    value := "concurrentValue"

    done := make(chan bool)

    go func() {
        cache.Set(key, value, 0)
        done <- true
    }()

    time.Sleep(1 * time.Millisecond) // Give some time for the Set operation

    <-done // Wait for the set operation to complete

    got, err := cache.Get(key)
    if err != nil {
        t.Fatalf("Concurrent Get() unexpected error: %v", err)
    }
    if got != value {
        t.Errorf("Concurrent Get() = %v, want %v", got, value)
    }
}

// TestCache_ReplaceExistingItem_Success tests replacing an existing item.
func TestCache_ReplaceExistingItem_Success(t *testing.T) {
    cache := NewCache(100*time.Millisecond, 0)
    defer cache.StopJanitor()

    key := "key"
    firstValue := "firstValue"
    secondValue := "secondValue"
    cache.Set(key, firstValue, 0)
    cache.Set(key, secondValue, 0) // Replace

    got, err := cache.Get(key)
    if err != nil {
        t.Fatalf("Get() after replace unexpected error: %v", err)
    }
    if got != secondValue {
        t.Errorf("Get() after replace = %v, want %v", got, secondValue)
    }
}

// TestCache_DeleteNonexistent_Success tests deleting a nonexistent item.
func TestCache_DeleteNonexistent_Success(t *testing.T) {
    cache := NewCache(100*time.Millisecond, 0)
    defer cache.StopJanitor()

    nonexistentKey := "nonexistentKey"
    cache.Delete(nonexistentKey) // Should not panic or error

    _, err := cache.Get(nonexistentKey)
    if !errors.Is(err, ErrItemNotFound) {
        t.Errorf("Delete() nonexistent item, expected %v, got %v", ErrItemNotFound, err)
    }
}

// TestCache_MultipleOperations_Success tests multiple operations performed on the cache.
func TestCache_MultipleOperations_Success(t *testing.T) {
    cache := NewCache(100*time.Millisecond, 0)
    defer cache.StopJanitor()

    key1 := "key1"
    value1 := "value1"
    key2 := "key2"
    value2 := "value2"

    cache.Set(key1, value1, 0)
    cache.Set(key2, value2, 0)

    cache.Delete(key1)

    _, err := cache.Get(key1)
    if !errors.Is(err, ErrItemNotFound) {
        t.Errorf("Expected %v error for key1, got %v", ErrItemNotFound, err)
    }

    got, err := cache.Get(key2)
    if err != nil {
        t.Fatalf("Unexpected error for key2: %v", err)
    }
    if got != value2 {
        t.Errorf("Expected value2 for key2, got %v", got)
    }
}

// TestCache_Clear_Success tests clearing the cache.
func TestCache_Clear_Success(t *testing.T) {
    cache := NewCache(100*time.Millisecond, 0)
    defer cache.StopJanitor()

    cache.Set("key1", "value1", 0)
    cache.Set("key2", "value2", 0)

    cache.Clear()

    if len(cache.items) != 0 {
        t.Errorf("Clear() failed, expected cache to be empty, got %d items", len(cache.items))
    }
}