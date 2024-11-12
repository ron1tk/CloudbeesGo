package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// mockStore simulates in-memory data storage for testing purposes
type mockStore struct {
	*InMemoryStore
}

func newMockStore() *mockStore {
	return &mockStore{
		InMemoryStore: &InMemoryStore{
			Users:      make(map[int]User),
			Tasks:      make(map[int]Task),
			tokenStore: make(map[string]int),
		},
	}
}

func (s *mockStore) SetupMockData() {
	// Adding a user
	s.Users[1] = User{
		ID:           1,
		Username:     "testuser",
		PasswordHash: HashPassword("testpass"),
		Token:        "validToken",
	}

	// Adding a task
	s.Tasks[1] = Task{
		ID:          1,
		Title:       "Test Task",
		Description: "Test Description",
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		OwnerID:     1,
	}

	// Linking token to user
	s.tokenStore["validToken"] = 1
}

func TestRegisterHandler(t *testing.T) {
	store := newMockStore()
	hc := &HandlerContext{Store: store.InMemoryStore}

	server := httptest.NewServer(http.HandlerFunc(hc.RegisterHandler))
	defer server.Close()

	tests := []struct {
		name       string
		username   string
		password   string
		wantStatus int
		wantMsg    string
	}{
		{"Valid Registration", "newuser", "newpass", http.StatusCreated, "User registered successfully"},
		{"Empty Username", "", "pass", http.StatusBadRequest, "Username and password are required"},
		{"Empty Password", "username", "", http.StatusBadRequest, "Username and password are required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(map[string]string{
				"username": tt.username,
				"password": tt.password,
			})
			response, body := testRequest(t, server.URL, "POST", bytes.NewBuffer(requestBody))
			defer response.Body.Close()

			if response.StatusCode != tt.wantStatus {
				t.Errorf("got status %d, want %d", response.StatusCode, tt.wantStatus)
			}

			var got Response
			if err := json.Unmarshal(body, &got); err != nil {
				t.Fatalf("unable to parse response: %v", err)
			}

			if !strings.Contains(got.Message, tt.wantMsg) {
				t.Errorf("expected message to contain %q, got %q", tt.wantMsg, got.Message)
			}
		})
	}
}

func TestLoginHandler(t *testing.T) {
	store := newMockStore()
	store.SetupMockData()
	hc := &HandlerContext{Store: store.InMemoryStore}

	server := httptest.NewServer(http.HandlerFunc(hc.LoginHandler))
	defer server.Close()

	tests := []struct {
		name       string
		username   string
		password   string
		wantStatus int
		wantMsg    string
	}{
		{"Valid Login", "testuser", "testpass", http.StatusOK, "Logged in successfully"},
		{"Invalid Password", "testuser", "wrongpass", http.StatusUnauthorized, "Invalid username or password"},
		{"Invalid Username", "wronguser", "testpass", http.StatusUnauthorized, "Invalid username or password"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(map[string]string{
				"username": tt.username,
				"password": tt.password,
			})
			response, body := testRequest(t, server.URL, "POST", bytes.NewBuffer(requestBody))
			defer response.Body.Close()

			if response.StatusCode != tt.wantStatus {
				t.Errorf("got status %d, want %d", response.StatusCode, tt.wantStatus)
			}

			var got Response
			if err := json.Unmarshal(body, &got); err != nil {
				t.Fatalf("unable to parse response: %v", err)
			}

			if !strings.Contains(got.Message, tt.wantMsg) {
				t.Errorf("expected message to contain %q, got %q", tt.wantMsg, got.Message)
			}
		})
	}
}

func testRequest(t *testing.T, url, method string, body *bytes.Buffer) (*http.Response, []byte) {
	t.Helper()

	client := &http.Client{}
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("could not read response body: %v", err)
	}

	return response, responseBody
}

// Note: The test suite is simplified for demonstration purposes.
// Real-world applications should include more detailed cases, including testing for concurrency issues,
// more complex scenarios, and integration tests that involve database operations or external services.
```