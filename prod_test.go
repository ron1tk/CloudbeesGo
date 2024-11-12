package main

import (
	"reflect"
	"testing"
	"time"
)

func TestCache_Set_Get_Delete(t *testing.T) {
	c := setupCache(1 * time.Minute)
	defer teardownCache(c)

	type args struct {
		key   string
		value interface{}
		ttl   time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Set and Get string value", args{"key1", "value1", 1 * time.Hour}, false},
		{"Set and Get integer value", args{"key2", 12345, 1 * time.Hour}, false},
		{"Set and Get struct value", args{"key3", struct{ Name string }{"John"}, 1 * time.Hour}, false},
		{"Set with negative ttl", args{"key4", "value4", -1 * time.Hour}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.Set(tt.args.key, tt.args.value, tt.args.ttl)
			if (err != nil) != tt.wantErr {
				t.Errorf("Cache.Set() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, err := c.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Cache.Get() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.args.value) {
				t.Errorf("Cache.Get() = %v, want %v", got, tt.args.value)
			}
		})
	}

	// Test Delete separately to ensure setup is correct
	t.Run("Delete existing key", func(t *testing.T) {
		keyToDelete := "keyToDelete"
		err := c.Set(keyToDelete, "valueToDelete", 1*time.Hour)
		if err != nil {
			t.Fatalf("Setup failed for delete test: %v", err)
		}

		c.Delete(keyToDelete)
		_, err = c.Get(keyToDelete)
		if err == nil {
			t.Errorf("Expected error after deleting item, got nil")
		}
	})
}

func TestCache_DeleteExpired(t *testing.T) {
	c := setupCache(1 * time.Millisecond)
	defer teardownCache(c)

	c.Set("keyToDeleteExpired", "value", 1*time.Nanosecond)
	time.Sleep(2 * time.Millisecond) // Ensure item has expired

	c.DeleteExpired()

	if _, found := c.items["keyToDeleteExpired"]; found {
		t.Errorf("DeleteExpired() did not delete the expired item.")
	}
}

func TestCache_Expiry(t *testing.T) {
	c := setupCache(10 * time.Millisecond)
	defer teardownCache(c)

	c.Set("keyExpireSoon", "value", 5*time.Millisecond)
	time.Sleep(6 * time.Millisecond) // Wait for item to expire

	_, err := c.Get("keyExpireSoon")
	if err == nil {
		t.Errorf("Expected error for expired item, got nil")
	}
}

func TestNewCache(t *testing.T) {
	c := NewCache(1 * time.Minute)
	if c == nil {
		t.Errorf("NewCache() failed to create a cache instance")
	}
}

func TestCache_StopJanitor(t *testing.T) {
	c := setupCache(1 * time.Millisecond)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("StopJanitor() caused panic: %v", r)
		}
	}()
	c.StopJanitor()
}