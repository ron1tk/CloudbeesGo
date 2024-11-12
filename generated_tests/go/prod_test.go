package main

import (
	"testing"
	"time"
)

func TestCache_Set_Get_SuccessCases(t *testing.T) {
	c := setupCache(10 * time.Minute)
	defer teardownCache(c)

	tests := []struct {
		name  string
		key   string
		value interface{}
		ttl   time.Duration
	}{
		{"Set and Get string value", "key1", "value1", 10 * time.Minute},
		{"Set and Get integer value", "key2", 12345, 10 * time.Minute},
		{"Set and Get struct value", "key3", struct{ Name string }{"John"}, 10 * time.Minute},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.Set(tt.key, tt.value, tt.ttl)
			if err != nil {
				t.Errorf("Cache.Set() error = %v, wantErr %v", err, false)
			}

			got, err := c.Get(tt.key)
			if err != nil {
				t.Errorf("Cache.Get() error = %v, wantErr %v", err, false)
			}
			if !reflect.DeepEqual(got, tt.value) {
				t.Errorf("Cache.Get() = %v, want %v", got, tt.value)
			}
		})
	}
}

func TestCache_Set_NegativeTTL(t *testing.T) {
	c := setupCache(10 * time.Minute)
	defer teardownCache(c)

	err := c.Set("keyNegativeTTL", "value", -1*time.Minute)
	if err == nil {
		t.Error("Cache.Set() with negative TTL did not return error, wantErr true")
	}
}

func TestCache_Get_NonExistentKey(t *testing.T) {
	c := setupCache(10 * time.Minute)
	defer teardownCache(c)

	_, err := c.Get("nonExistentKey")
	if err == nil {
		t.Error("Cache.Get() for non-existent key did not return error, wantErr true")
	}
}

func TestCache_Delete(t *testing.T) {
	c := setupCache(10 * time.Minute)
	defer teardownCache(c)

	key := "keyToDelete"
	err := c.Set(key, "valueToDelete", 10*time.Minute)
	if err != nil {
		t.Fatalf("Failed to set up for delete test: %v", err)
	}

	c.Delete(key)
	_, err = c.Get(key)
	if err == nil {
		t.Error("Cache.Get() after delete did not return error, wantErr true")
	}
}

func TestCache_DeleteExpired_Success(t *testing.T) {
	c := setupCache(1 * time.Millisecond)
	defer teardownCache(c)

	c.Set("keyToDeleteExpired", "value", 1*time.Nanosecond)
	time.Sleep(2 * time.Millisecond) // Ensure item has expired

	c.DeleteExpired()

	_, found := c.items["keyToDeleteExpired"]
	if found {
		t.Error("DeleteExpired() did not delete the expired item.")
	}
}

func TestCache_Expiry(t *testing.T) {
	c := setupCache(10 * time.Millisecond)
	defer teardownCache(c)

	c.Set("keyExpireSoon", "value", 5*time.Millisecond)
	time.Sleep(6 * time.Millisecond) // Wait for item to expire

	_, err := c.Get("keyExpireSoon")
	if err == nil {
		t.Error("Expected error for expired item, got nil")
	}
}

func TestNewCache(t *testing.T) {
	c := NewCache(10 * time.Minute)
	if c == nil {
		t.Error("NewCache() failed to create a cache instance")
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