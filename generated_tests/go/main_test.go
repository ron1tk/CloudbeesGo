// main_test.go
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Mock DB
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) First(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(out, where)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Save(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Delete(value interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(value, where)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Model(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	call := m.Called(query, args)
	return call.Get(0).(*gorm.DB)
}

func (m *MockDB) Limit(limit int) *gorm.DB {
	args := m.Called(limit)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Offset(offset int) *gorm.DB {
	args := m.Called(offset)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(out, where)
	return args.Get(0).(*gorm.DB)
}

func setup() *App {
	app := &App{Router: mux.NewRouter(), DB: new(MockDB), Validator: validator.New(), JWTSecret: "testsecret", TokenExpiry: time.Minute * 60}
	app.initializeRoutes()
	return app
}

func TestHealthCheckHandler(t *testing.T) {
	app := setup()
	req, _ := http.NewRequest("GET", "/health", nil)
	response := executeRequest(req, app)

	checkResponseCode(t, http.StatusOK, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	assert.Equal(t, "healthy", m["status"])
}

func TestLoginHandler_Success(t *testing.T) {
	app := setup()

	creds := Credentials{
		Username: "admin",
		Password: "password",
	}

	body, _ := json.Marshal(creds)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	response := executeRequest(req, app)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestLoginHandler_Failure(t *testing.T) {
	app := setup()

	creds := Credentials{
		Username: "wrong",
		Password: "creds",
	}

	body, _ := json.Marshal(creds)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	response := executeRequest(req, app)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

func TestCreateUserHandler_Success(t *testing.T) {
	app := setup()

	mockDB := app.DB.(*MockDB)
	mockDB.On("Create", mock.AnythingOfType("*main.User")).Return(&gorm.DB{Error: nil})

	userInput := UserInput{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	body, _ := json.Marshal(userInput)
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", generateTestToken(app)))
	response := executeRequest(req, app)

	checkResponseCode(t, http.StatusCreated, response.Code)
}

func TestCreateUserHandler_Failure(t *testing.T) {
	app := setup()

	userInput := UserInput{
		FirstName: "",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	body, _ := json.Marshal(userInput)
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", generateTestToken(app)))
	response := executeRequest(req, app)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestGetUserHandler_NotFound(t *testing.T) {
	app := setup()

	mockDB := app.DB.(*MockDB)
	mockDB.On("First", mock.Anything, mock.Anything).Return(&gorm.DB{Error: gorm.ErrRecordNotFound})

	req, _ := http.NewRequest("GET", "/api/users/123", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", generateTestToken(app)))
	response := executeRequest(req, app)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func executeRequest(req *http.Request, app *App) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	app.Router.ServeHTTP(recorder, req)
	return recorder
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func generateTestToken(app *App) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "admin",
		"exp":      time.Now().Add(time.Minute * 5).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(app.JWTSecret))
	return tokenString
}

func TestMain(m *testing.M) {
	// Setup test environment
	cfg = Config{
		ServerPort:      "8080",
		JWTSecret:       "testsecret",
		RateLimit:       "5-S",
		DatabasePath:    ":memory:",
		LogLevel:        "panic",
		TokenExpiryMins: 60,
	}

	app := App{}
	err := app.Initialize()
	if err != nil {
		fmt.Println("Failed to initialize the test application.")
		return
	}

	code := m.Run()

	os.Exit(code)
}