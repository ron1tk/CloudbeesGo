// main_test.go
package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	main "path/to/your/main/package"
)

var (
	app *main.App
)

func setup() {
	app = &main.App{}
	err := app.Initialize()
	if err != nil {
		fmt.Printf("Failed to initialize the app: %v\n", err)
	}
}

func tearDown() {
	// Any cleanup logic goes here.
}

func TestHealthCheckHandler(t *testing.T) {
	setup()
	defer tearDown()

	req, _ := http.NewRequest("GET", "/health", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	assert.Equal(t, "healthy", m["status"], "Expected health check status to be healthy.")
}

func TestVersionHandler(t *testing.T) {
	setup()
	defer tearDown()

	req, _ := http.NewRequest("GET", "/version", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	assert.Equal(t, "1.0.0", m["version"], "Expected version to be 1.0.0.")
}

func TestLoginHandler(t *testing.T) {
	setup()
	defer tearDown()

	var jsonStr = []byte(`{"username":"admin","password":"password"}`)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	tokenString, ok := m["token"]
	assert.True(t, ok, "Expected token to be returned")
	assert.NotEmpty(t, tokenString, "Expected non-empty token string")
}

func TestLoginHandlerInvalidCredentials(t *testing.T) {
	setup()
	defer tearDown()

	var jsonStr = []byte(`{"username":"wrong","password":"credentials"}`)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

func TestCreateUserHandler(t *testing.T) {
	setup()
	defer tearDown()

	var jsonStr = []byte(`{"first_name":"John","last_name":"Doe","email":"johndoe@example.com"}`)
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	// Mocking JWT Authentication
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "admin",
	})
	tokenString, _ := token.SignedString([]byte("mysecretkey"))
	req.Header.Set("Authorization", "Bearer "+tokenString)

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
}

func TestGetUserHandler(t *testing.T) {
	setup()
	defer tearDown()

	req, _ := http.NewRequest("GET", "/api/users/1", nil)

	// Mocking JWT Authentication
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "admin",
	})
	tokenString, _ := token.SignedString([]byte("mysecretkey"))
	req.Header.Set("Authorization", "Bearer "+tokenString)

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateUserHandler(t *testing.T) {
	// Similar setup to TestCreateUserHandler with a PUT request to "/api/users/{id}"
}

func TestDeleteUserHandler(t *testing.T) {
	// Similar setup to TestGetUserHandler with a DELETE request to "/api/users/{id}"
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	app.Router.ServeHTTP(recorder, req)
	return recorder
}

func checkResponseCode(t *testing.T, expected, actual int) {
	assert.Equal(t, expected, actual, fmt.Sprintf("Expected response code %d. Got %d\n", expected, actual))
}