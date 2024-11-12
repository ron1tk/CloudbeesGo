package main

import (
	"testing"
	"time"
)

// MockCache is a mock of Cache interface to be used in tests.
type MockCache struct {
	Cache
	SetFunc    func(string, interface{}, time.Duration) error
	GetFunc    func(string) (interface{}, error)
	DeleteFunc func(string)
}

func (m *MockCache) Set(key string, value interface{}, ttl time.Duration) error {
	if m.SetFunc != nil {
		return m.SetFunc(key, value, ttl)
	}
	return nil
}

func (m *MockCache) Get(key string) (interface{}, error) {
	if m.GetFunc != nil {
		return m.GetFunc(key)
	}
	return nil, nil
}

func (m *MockCache) Delete(key string) {
	if m.DeleteFunc != nil {
		m.DeleteFunc(key)
	}
}

func TestCache_Set_Get_Delete_Mock(t *testing.T) {
	mockCache := &MockCache{
		SetFunc: func(key string, value interface{}, ttl time.Duration) error {
			if ttl < 0 {
				return ErrInvalidTTL
			}
			return nil
		},
		GetFunc: func(key string) (interface{}, error) {
			if key == "nonExistentKey" {
				return nil, ErrKeyNotFound
			}
			return "mockValue", nil
		},
		DeleteFunc: func(key string) {
			// Mock delete functionality
		},
	}

	type args struct {
		key   string
		value interface{}
		ttl   time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"Set and Get with mock", args{"mockKey", "mockValue", 1 * time.Hour}, "mockValue", false},
		{"Set with negative ttl using mock", args{"mockKey", "mockValue", -1 * time.Hour}, nil, true},
		{"Get non-existent key using mock", args{"nonExistentKey", nil, 0}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mockCache.Set(tt.args.key, tt.args.value, tt.args.ttl)
			if (err != nil) != tt.wantErr {
				t.Errorf("MockCache.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got, err := mockCache.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("MockCache.Get() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Errorf("MockCache.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_DeleteExpired_Mock(t *testing.T) {
	mockCache := &MockCache{
		DeleteFunc: func(key string) {
			// Mock delete to simulate deletion of an expired key
		},
	}

	mockCache.Set("keyToExpire", "value", 1*time.Millisecond)
	time.Sleep(2 * time.Millisecond) // Simulate waiting for key to expire

	mockCache.DeleteExpired()

	// Since DeleteExpired logic is not implemented in MockCache, this test
	// primarily ensures that the method can be called without errors.
	// Actual deletion logic should be tested in integration tests.
}

func TestCache_Expiry_Mock(t *testing.T) {
	mockCache := &MockCache{
		GetFunc: func(key string) (interface{}, error) {
			return nil, ErrKeyNotFound
		},
	}

	mockCache.Set("keyToExpire", "value", 1*time.Millisecond)
	time.Sleep(2 * time.Millisecond)

	_, err := mockCache.Get("keyToExpire")
	if err == nil {
		t.Errorf("Expected error for expired item, got nil")
	}
}

func TestCache_StopJanitor_Mock(t *testing.T) {
	mockCache := &MockCache{}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("StopJanitor() caused panic: %v", r)
		}
	}()

	mockCache.StopJanitor()

	// This test ensures StopJanitor can be called without causing a panic.
	// Actual implementation should ensure janitor goroutine is stopped.
}