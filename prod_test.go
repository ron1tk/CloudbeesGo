// prod_test.go
package cache

import (
    "errors"
    "reflect"
    "testing"
    "time"
)

func TestCache_SetWithExpiration(t *testing.T) {
    // Provide both cleanupInterval and defaultDuration
    cache := NewCache(10*time.Millisecond, 0)
    defer cache.StopJanitor()

    key := "expiringKey"
    value := "expiringValue"
    expiration := 5 * time.Millisecond
    cache.Set(key, value, expiration)

    time.Sleep(6 * time.Millisecond) // Wait for the item to expire

    _, err := cache.Get(key)
    if !errors.Is(err, ErrItemExpired) {
        t.Errorf("Expected error %v for expired item, got %v", ErrItemExpired, err)
    }
}

func TestCache_SetAndGetConcurrently(t *testing.T) {
    cache := NewCache(100*time.Millisecond, 0)
    defer cache.StopJanitor()

    key := "concurrentKey"
    value := "concurrentValue"

    go cache.Set(key, value, 0)

    time.Sleep(1 * time.Millisecond) // Give some time for the Set operation

    got, err := cache.Get(key)
    if err != nil {
        t.Errorf("Concurrent Get() unexpected error: %v", err)
    }
    if !reflect.DeepEqual(got, value) {
        t.Errorf("Concurrent Get() = %v, want %v", got, value)
    }
}

func TestCache_ReplaceExistingItem(t *testing.T) {
    cache := NewCache(100*time.Millisecond, 0)
    defer cache.StopJanitor()

    key := "key"
    firstValue := "firstValue"
    secondValue := "secondValue"
    cache.Set(key, firstValue, 0)
    cache.Set(key, secondValue, 0) // Replace

    got, err := cache.Get(key)
    if err != nil {
        t.Errorf("Get() after replace unexpected error: %v", err)
    }
    if !reflect.DeepEqual(got, secondValue) {
        t.Errorf("Get() after replace = %v, want %v", got, secondValue)
    }
}

func TestCache_DeleteNonexistent(t *testing.T) {
    cache := NewCache(100*time.Millisecond, 0)
    defer cache.StopJanitor()

    nonexistentKey := "nonexistentKey"
    cache.Delete(nonexistentKey) // Should not panic or error

    _, err := cache.Get(nonexistentKey)
    if !errors.Is(err, ErrItemNotFound) {
        t.Errorf("Delete() nonexistent item, expected %v, got %v", ErrItemNotFound, err)
    }
}

func TestCache_MultipleOperations(t *testing.T) {
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
        t.Errorf("Unexpected error for key2: %v", err)
    }
    if got != value2 {
        t.Errorf("Expected value2 for key2, got %v", got)
    }
}

func TestCache_Clear(t *testing.T) {
    cache := NewCache(100*time.Millisecond, 0)
    defer cache.StopJanitor()

    cache.Set("key1", "value1", 0)
    cache.Set("key2", "value2", 0)

    cache.Clear()

    if len(cache.items) != 0 {
        t.Errorf("Clear() failed, expected cache to be empty, got %d items", len(cache.items))
    }
}
