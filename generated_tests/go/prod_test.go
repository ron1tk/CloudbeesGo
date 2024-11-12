package cache

import (
    "reflect"
    "testing"
    "time"
)

func TestNewCache(t *testing.T) {
    cleanupInterval := 10 * time.Millisecond
    cache := NewCache(cleanupInterval)

    if cache == nil {
        t.Error("Expected NewCache to create a non-nil cache instance")
    }
}

func TestCache_SetAndGet(t *testing.T) {
    cache := NewCache(10 * time.Millisecond)
    defer cache.StopJanitor()

    tests := []struct {
        name      string
        key       string
        value     interface{}
        duration  time.Duration
        want      interface{}
        wantError bool
    }{
        {"Set and Get existing item", "key1", "value1", 0, "value1", false},
        {"Get non-existing item", "key2", nil, 0, nil, true},
        {"Set and Get expired item", "key3", "value3", 1 * time.Millisecond, nil, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cache.Set(tt.key, tt.value, tt.duration)
            time.Sleep(2 * time.Millisecond) // Ensure some items expire

            got, err := cache.Get(tt.key)
            if (err != nil) != tt.wantError {
                t.Errorf("Cache.Get() error = %v, wantError %v", err, tt.wantError)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("Cache.Get() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestCache_Delete(t *testing.T) {
    cache := NewCache(10 * time.Millisecond)
    defer cache.StopJanitor()

    key := "key1"
    value := "value1"
    cache.Set(key, value, 0)

    cache.Delete(key)

    _, err := cache.Get(key)
    if err == nil {
        t.Errorf("Cache.Delete() failed, expected error but got nil")
    }
}

func TestCache_DeleteExpired(t *testing.T) {
    cache := NewCache(1 * time.Millisecond)
    defer cache.StopJanitor()

    cache.Set("key1", "value1", 1*time.Millisecond)

    time.Sleep(2 * time.Millisecond) // Ensure the item expires

    cache.DeleteExpired()

    _, err := cache.Get("key1")
    if err == nil {
        t.Errorf("Cache.DeleteExpired() failed, expected error but got nil")
    }
}

func TestStopJanitor(t *testing.T) {
    cache := NewCache(1 * time.Millisecond)

    // Difficult to test the stopping without introducing some kind of synchronization
    // mechanism to observe the goroutine stopping. This would likely require
    // modifications to the prod code for testing purposes, which isn't ideal.
    // For now, this test just ensures calling StopJanitor doesn't result in panic.
    cache.StopJanitor()
}