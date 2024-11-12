// main_test.go
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestHealthCheckHandler(t *testing.T) {
	app := setup()
	req, _ := http.NewRequest("GET", "/health", nil)
	response := executeRequest(req, app)

	checkResponseCode(t, http.StatusOK, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	assert.Equal(t, "healthy", m["status"], "Expected health check status to be 'healthy'")
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
	mockDB.AssertExpectations(t)
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
	mockDB.AssertExpectations(t)
}

// Additional Test Cases

func TestUpdateUserHandler_Success(t *testing.T) {
	app := setup()

	mockDB := app.DB.(*MockDB)
	mockDB.On("Model", mock.AnythingOfType("*main.User")).Return(&gorm.DB{}).Once()
	mockDB.On("Where", "id = ?", mock.Anything).Return(&gorm.DB{}).Once()
	mockDB.On("Save", mock.AnythingOfType("*main.User")).Return(&gorm.DB{Error: nil})

	userUpdate := UserInput{
		FirstName: "UpdatedName",
	}

	body, _ := json.Marshal(userUpdate)
	req, _ := http.NewRequest("PUT", "/api/users/123", bytes.NewBuffer(body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", generateTestToken(app)))
	response := executeRequest(req, app)

	checkResponseCode(t, http.StatusOK, response.Code)
	mockDB.AssertExpectations(t)
}

func TestUpdateUserHandler_Failure_NotFound(t *testing.T) {
	app := setup()

	mockDB := app.DB.(*MockDB)
	mockDB.On("Model", mock.AnythingOfType("*main.User")).Return(&gorm.DB{}).Once()
	mockDB.On("Where", "id = ?", mock.Anything).Return(&gorm.DB{Error: gorm.ErrRecordNotFound})

	userUpdate := UserInput{
		FirstName: "NonExistent",
	}

	body, _ := json.Marshal(userUpdate)
	req, _ := http.NewRequest("PUT", "/api/users/999", bytes.NewBuffer(body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", generateTestToken(app)))
	response := executeRequest(req, app)

	checkResponseCode(t, http.StatusNotFound, response.Code)
	mockDB.AssertExpectations(t)
}

func TestDeleteUserHandler_Success(t *testing.T) {
	app := setup()

	mockDB := app.DB.(*MockDB)
	mockDB.On("Delete", mock.AnythingOfType("*main.User"), mock.Anything).Return(&gorm.DB{Error: nil})

	req, _ := http.NewRequest("DELETE", "/api/users/123", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", generateTestToken(app)))
	response := executeRequest(req, app)

	checkResponseCode(t, http.StatusOK, response.Code)
	mockDB.AssertExpectations(t)
}

func TestDeleteUserHandler_Failure_NotFound(t *testing.T) {
	app := setup()

	mockDB := app.DB.(*MockDB)
	mockDB.On("Delete", mock.AnythingOfType("*main.User"), mock.Anything).Return(&gorm.DB{Error: gorm.ErrRecordNotFound})

	req, _ := http.NewRequest("DELETE", "/api/users/999", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", generateTestToken(app)))
	response := executeRequest(req, app)

	checkResponseCode(t, http.StatusNotFound, response.Code)
	mockDB.AssertExpectations(t)
}