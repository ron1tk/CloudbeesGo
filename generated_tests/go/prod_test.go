package cache

import (
    "reflect"
    "testing"
    "time"
)

// TestCache_SetAndGet_Success checks if setting and getting a cache item works as expected.
func TestCache_SetAndGet_Success(t *testing.T) {
    c := setupCache(1 * time.Minute)
    defer teardownCache(c)

    testKey := "key1"
    testValue := "value1"
    c.Set(testKey, testValue, 1*time.Hour)

    gotValue, err := c.Get(testKey)
    if err != nil {
        t.Fatalf("Get() error = %v, wantErr %v", err, false)
    }
    if !reflect.DeepEqual(gotValue, testValue) {
        t.Errorf("Get() = %v, want %v", gotValue, testValue)
    }
}

// TestCache_Get_ItemNotFound verifies that getting a non-existent cache item returns an error.
func TestCache_Get_ItemNotFound(t *testing.T) {
    c := setupCache(1 * time.Minute)
    defer teardownCache(c)

    _, err := c.Get("nonexistent")
    if err == nil {
        t.Error("Expected error for getting non-existent item, got nil")
    }
}

// TestCache_Get_ItemExpired checks behavior when trying to get an expired cache item.
func TestCache_Get_ItemExpired(t *testing.T) {
    c := setupCache(1 * time.Millisecond)
    defer teardownCache(c)

    c.Set("key2", "value2", 1*time.Nanosecond)
    time.Sleep(2 * time.Millisecond) // Wait for item to expire

    _, err := c.Get("key2")
    if err == nil {
        t.Error("Expected error for getting expired item, got nil")
    }
}

// TestCache_Delete verifies that deleting an item works correctly.
func TestCache_Delete(t *testing.T) {
    c := setupCache(1 * time.Minute)
    defer teardownCache(c)

    c.Set("key3", "value3", 1*time.Hour)
    c.Delete("key3")

    _, err := c.Get("key3")
    if err == nil {
        t.Error("Expected error after deleting item, got nil")
    }
}

// TestCache_DeleteExpired ensures expired items are properly deleted from the cache.
func TestCache_DeleteExpired(t *testing.T) {
    c := setupCache(1 * time.Millisecond)
    defer teardownCache(c)

    c.Set("key4", "value4", 1*time.Nanosecond)
    time.Sleep(2 * time.Millisecond) // Wait for item to expire

    c.DeleteExpired()

    _, found := c.items["key4"]
    if found {
        t.Error("Expected expired item to be deleted, but it was found")
    }
}

// TestNewCache checks that a new cache is created without errors.
func TestNewCache(t *testing.T) {
    c := setupCache(1 * time.Minute)
    defer teardownCache(c)

    if c == nil {
        t.Error("NewCache() returned nil, expected a valid Cache instance")
    }
}

// TestStopJanitor verifies that stopping the janitor does not result in a panic or error.
func TestStopJanitor(t *testing.T) {
    c := setupCache(1 * time.Millisecond)

    defer func() {
        if r := recover(); r != nil {
            t.Error("Expected StopJanitor to not panic, but it did")
        }
    }()
    c.StopJanitor()
}

// TestCache_Set_OverwriteExisting ensures that setting a key that already exists overwrites the existing value.
func TestCache_Set_OverwriteExisting(t *testing.T) {
    c := setupCache(1 * time.Minute)
    defer teardownCache(c)

    key := "overwriteKey"
    firstValue := "firstValue"
    secondValue := "secondValue"

    c.Set(key, firstValue, 1*time.Hour)
    c.Set(key, secondValue, 1*time.Hour)

    gotValue, err := c.Get(key)
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
    if gotValue != secondValue {
        t.Errorf("Expected value to be overwritten to %v, got %v", secondValue, gotValue)
    }
}

// TestCache_Concurrency checks if the cache behaves correctly under concurrent access.
func TestCache_Concurrency(t *testing.T) {
    c := setupCache(1 * time.Minute)
    defer teardownCache(c)

    key := "concurrentKey"
    value := "concurrentValue"

    go func() {
        c.Set(key, value, 1*time.Hour)
    }()

    go func() {
        c.Delete(key)
    }()

    time.Sleep(100 * time.Millisecond) // Give goroutines time to execute

    _, err := c.Get(key)
    if err == nil {
        t.Error("Expected error for getting a possibly deleted item, got nil")
    }
}