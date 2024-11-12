package cache

import (
    "reflect"
    "testing"
    "time"
)

func setupCache(defaultExpiration time.Duration) *Cache {
    // This setup function should be implemented to initialize the cache with a default expiration time.
    return NewCache(defaultExpiration)
}

func teardownCache(c *Cache) {
    // This teardown function can be used to clean up resources or run any necessary tear down steps.
    c.StopJanitor()
}

func TestCache_Set_Success(t *testing.T) {
    c := setupCache(1 * time.Minute)
    defer teardownCache(c)

    err := c.Set("key", "value", 1*time.Hour)
    if err != nil {
        t.Errorf("Set() error = %v, wantErr %v", err, false)
    }
}

func TestCache_Set_Failure(t *testing.T) {
    c := setupCache(1 * time.Minute)
    defer teardownCache(c)

    // Assuming Set can fail under certain conditions, which should be mocked or simulated if possible.
}

func TestCache_Get_Success(t *testing.T) {
    c := setupCache(1 * time.Minute)
    defer teardownCache(c)

    testKey := "key"
    testValue := "value"
    c.Set(testKey, testValue, 1*time.Hour)

    got, err := c.Get(testKey)
    if err != nil {
        t.Fatalf("Get() error = %v, wantErr %v", err, false)
    }
    if !reflect.DeepEqual(got, testValue) {
        t.Errorf("Get() got = %v, want %v", got, testValue)
    }
}

func TestCache_Get_Failure(t *testing.T) {
    c := setupCache(1 * time.Minute)
    defer teardownCache(c)

    _, err := c.Get("unknown")
    if err == nil {
        t.Errorf("Get() expected error, got nil")
    }
}

func TestCache_Delete_Success(t *testing.T) {
    c := setupCache(1 * time.Minute)
    defer teardownCache(c)

    testKey := "keyToDelete"
    c.Set(testKey, "value", 1*time.Hour)
    c.Delete(testKey)

    _, err := c.Get(testKey)
    if err == nil {
        t.Errorf("Delete() expected error, got nil")
    }
}

func TestCache_DeleteExpired(t *testing.T) {
    c := setupCache(10 * time.Millisecond)
    defer teardownCache(c)

    testKey := "expiredKey"
    c.Set(testKey, "value", 1*time.Millisecond)

    time.Sleep(20 * time.Millisecond) // ensure the item is expired
    c.DeleteExpired()

    _, err := c.Get(testKey)
    if err == nil {
        t.Errorf("DeleteExpired() expected error, got nil")
    }
}

func TestCache_Clear(t *testing.T) {
    c := setupCache(1 * time.Minute)
    defer teardownCache(c)

    c.Set("key1", "value1", 1*time.Hour)
    c.Set("key2", "value2", 1*time.Hour)
    c.Clear()

    if len(c.items) != 0 {
        t.Errorf("Clear() did not clear the cache, items count = %d", len(c.items))
    }
}

func TestCache_ItemCount(t *testing.T) {
    c := setupCache(1 * time.Minute)
    defer teardownCache(c)

    c.Set("key1", "value1", 1*time.Hour)
    c.Set("key2", "value2", 1*time.Hour)

    count := c.ItemCount()
    if count != 2 {
        t.Errorf("ItemCount() = %d, want %d", count, 2)
    }
}