package cache

import (
    "reflect"
    "testing"
    "time"
)

func TestCache_SetAndGet_Success(t *testing.T) {
    c := setupCache(1 * time.Minute)
    defer teardownCache(c)

    tests := []struct {
        name  string
        key   string
        value interface{}
        ttl   time.Duration
    }{
        {"String value", "key1", "value1", 1 * time.Hour},
        {"Integer value", "key2", 12345, 1 * time.Hour},
        {"Struct value", "key3", struct{ Name string }{"John"}, 1 * time.Hour},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c.Set(tt.key, tt.value, tt.ttl)

            got, err := c.Get(tt.key)
            if err != nil {
                t.Errorf("Get() error = %v, wantErr %v", err, false)
            }
            if !reflect.DeepEqual(got, tt.value) {
                t.Errorf("Get() = %v, want %v", got, tt.value)
            }
        })
    }
}

func TestCache_Get_ItemNotFound_Error(t *testing.T) {
    c := setupCache(1 * time.Minute)
    defer teardownCache(c)

    _, err := c.Get("nonexistent")
    if err == nil {
        t.Errorf("Expected error for non-existing item, got nil")
    }
}

func TestCache_SetAndGet_ItemExpired_Error(t *testing.T) {
    c := setupCache(1 * time.Millisecond)
    defer teardownCache(c)

    c.Set("keyExpired", "value", 1*time.Nanosecond)

    time.Sleep(2 * time.Millisecond) // Wait for item to expire

    _, err := c.Get("keyExpired")
    if err == nil {
        t.Errorf("Expected error for expired item, got nil")
    }
}

func TestCache_Delete_Success(t *testing.T) {
    c := setupCache(1 * time.Minute)
    defer teardownCache(c)

    c.Set("keyToDelete", "value", 1*time.Hour)
    c.Delete("keyToDelete")

    _, err := c.Get("keyToDelete")
    if err == nil {
        t.Errorf("Expected error after deleting item, got nil")
    }
}

func TestCache_DeleteExpired_Success(t *testing.T) {
    c := setupCache(1 * time.Millisecond)
    defer teardownCache(c)

    c.Set("keyToDeleteExpired", "value", 1*time.Nanosecond)
    time.Sleep(2 * time.Millisecond) // Wait for item to expire

    c.DeleteExpired()

    if _, found := c.items["keyToDeleteExpired"]; found {
        t.Errorf("DeleteExpired() did not delete the expired item.")
    }
}

func TestNewCache_CreatesCacheInstance(t *testing.T) {
    c := NewCache(1 * time.Minute)
    if c == nil {
        t.Errorf("NewCache() failed to create a cache instance")
    }
}

func TestCache_StopJanitor_NoPanic(t *testing.T) {
    c := setupCache(1 * time.Millisecond)
    defer func() {
        if r := recover(); r != nil {
            t.Errorf("StopJanitor() caused panic: %v", r)
        }
    }()
    c.StopJanitor()
}