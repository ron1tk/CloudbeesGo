package cache

import (
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestCache_NewCache(t *testing.T) {
	cache := NewCache(100 * time.Millisecond)
	defer cache.StopJanitor()

	if cache == nil {
		t.Error("NewCache() failed, expected a new cache instance, got nil")
	}
}

func TestCache_SetWithExpiration_ExpiredItem(t *testing.T) {
	cache := NewCache(10 * time.Millisecond)
	defer cache.StopJanitor()

	cache.Set("expiringKey", "expiringValue", 5*time.Millisecond)

	time.Sleep(10 * time.Millisecond) // Wait for the item to definitely expire

	_, err := cache.Get("expiringKey")
	if !errors.Is(err, ErrItemNotFound) {
		t.Errorf("Expected ErrItemNotFound for expired item, got %v", err)
	}
}

func TestCache_SetWithExpiration_NonExpiredItem(t *testing.T) {
	cache := NewCache(20 * time.Millisecond)
	defer cache.StopJanitor()

	key := "key"
	value := "value"
	cache.Set(key, value, 15*time.Millisecond)

	time.Sleep(10 * time.Millisecond) // Item should still exist

	got, err := cache.Get(key)
	if err != nil {
		t.Fatalf("Get() unexpected error: %v", err)
	}
	if got != value {
		t.Errorf("Get() = %v, want %v", got, value)
	}
}

func TestCache_SetAndGetConcurrently_MultipleRoutines(t *testing.T) {
	cache := NewCache(100 * time.Millisecond)
	defer cache.StopJanitor()

	key := "concurrentKey"
	value := "concurrentValue"

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		cache.Set(key, value, 0)
	}()

	wg.Wait()

	time.Sleep(1 * time.Millisecond) // Ensure Set() completes in goroutine

	got, err := cache.Get(key)
	if err != nil {
		t.Errorf("Concurrent Get() unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, value) {
		t.Errorf("Concurrent Get() = %v, want %v", got, value)
	}
}

func TestCache_DeleteExisting(t *testing.T) {
	cache := NewCache(100 * time.Millisecond)
	defer cache.StopJanitor()

	key := "existingKey"
	value := "existingValue"
	cache.Set(key, value, 0)

	cache.Delete(key)

	_, err := cache.Get(key)
	if !errors.Is(err, ErrItemNotFound) {
		t.Errorf("Delete() existing item, expected %v, got %v", ErrItemNotFound, err)
	}
}

func TestCache_Clear_WithItems(t *testing.T) {
	cache := NewCache(100 * time.Millisecond)
	defer cache.StopJanitor()

	cache.Set("key1", "value1", 0)
	cache.Set("key2", "value2", 0)

	cache.Clear()

	if len(cache.items) != 0 {
		t.Errorf("Clear() failed, expected cache to be empty, got %d items", len(cache.items))
	}
}

func TestCache_Clear_EmptyCache(t *testing.T) {
	cache := NewCache(100 * time.Millisecond)
	defer cache.StopJanitor()

	cache.Clear() // Clearing an already empty cache

	if len(cache.items) != 0 {
		t.Errorf("Clear() on an empty cache, expected cache to be empty, got %d items", len(cache.items))
	}
}

func TestCache_Get_NonExisting(t *testing.T) {
	cache := NewCache(100 * time.Millisecond)
	defer cache.StopJanitor()

	_, err := cache.Get("nonExistingKey")
	if !errors.Is(err, ErrItemNotFound) {
		t.Errorf("Get() non-existing item, expected %v, got %v", ErrItemNotFound, err)
	}
}