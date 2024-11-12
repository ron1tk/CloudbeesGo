package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)


// Constants for token generation and expiration
const (
	TokenLength    = 32
	TokenExpiry    = time.Hour * 24
	ServerPort     = "8000"
	ContextUserKey = "user"
)

// User represents a user in the system.
type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Token        string `json:"token,omitempty"`
}

// Task represents a task with various attributes.
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	OwnerID     int       `json:"owner_id"`
}

// Response represents a standard API response.
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// InMemoryStore holds the users and tasks with mutexes for concurrency safety.
type InMemoryStore struct {
	Users      map[int]User
	Tasks      map[int]Task
	userMutex  sync.RWMutex
	taskMutex  sync.RWMutex
	userIDSeq  int
	taskIDSeq  int
	tokenStore map[string]int // token to user ID
	tokenMux   sync.RWMutex
}

// NewStore initializes and returns a new InMemoryStore.
func NewStore() *InMemoryStore {
	return &InMemoryStore{
		Users:      make(map[int]User),
		Tasks:      make(map[int]Task),
		tokenStore: make(map[string]int),
	}
}

// AddUser adds a new user to the store.
func (s *InMemoryStore) AddUser(username, password string) (User, error) {
	s.userMutex.Lock()
	defer s.userMutex.Unlock()

	// Check if username already exists
	for _, user := range s.Users {
		if user.Username == username {
			return User{}, errors.New("username already exists")
		}
	}

	// Hash the password
	hashedPassword := HashPassword(password)

	s.userIDSeq++
	user := User{
		ID:           s.userIDSeq,
		Username:     username,
		PasswordHash: hashedPassword,
	}
	s.Users[user.ID] = user
	return user, nil
}

// AuthenticateUser authenticates a user and returns the user if successful.
func (s *InMemoryStore) AuthenticateUser(username, password string) (User, error) {
	s.userMutex.RLock()
	defer s.userMutex.RUnlock()

	for _, user := range s.Users {
		if user.Username == username {
			if CheckPasswordHash(password, user.PasswordHash) {
				return user, nil
			}
			return User{}, errors.New("invalid password")
		}
	}
	return User{}, errors.New("user not found")
}

// AddTask adds a new task to the store.
func (s *InMemoryStore) AddTask(title, description string, ownerID int) Task {
	s.taskMutex.Lock()
	defer s.taskMutex.Unlock()

	s.taskIDSeq++
	now := time.Now()
	task := Task{
		ID:          s.taskIDSeq,
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
		OwnerID:     ownerID,
	}
	s.Tasks[task.ID] = task
	return task
}

// GetTask retrieves a task by ID.
func (s *InMemoryStore) GetTask(taskID, ownerID int) (Task, error) {
	s.taskMutex.RLock()
	defer s.taskMutex.RUnlock()

	task, exists := s.Tasks[taskID]
	if !exists {
		return Task{}, errors.New("task not found")
	}
	if task.OwnerID != ownerID {
		return Task{}, errors.New("unauthorized access to task")
	}
	return task, nil
}

// UpdateTask updates an existing task.
func (s *InMemoryStore) UpdateTask(taskID, ownerID int, updated Task) (Task, error) {
	s.taskMutex.Lock()
	defer s.taskMutex.Unlock()

	task, exists := s.Tasks[taskID]
	if !exists {
		return Task{}, errors.New("task not found")
	}
	if task.OwnerID != ownerID {
		return Task{}, errors.New("unauthorized access to task")
	}

	// Update fields if provided
	if updated.Title != "" {
		task.Title = updated.Title
	}
	if updated.Description != "" {
		task.Description = updated.Description
	}
	task.Completed = updated.Completed
	task.UpdatedAt = time.Now()

	s.Tasks[taskID] = task
	return task, nil
}

// DeleteTask removes a task from the store.
func (s *InMemoryStore) DeleteTask(taskID, ownerID int) error {
	s.taskMutex.Lock()
	defer s.taskMutex.Unlock()

	task, exists := s.Tasks[taskID]
	if !exists {
		return errors.New("task not found")
	}
	if task.OwnerID != ownerID {
		return errors.New("unauthorized access to task")
	}

	delete(s.Tasks, taskID)
	return nil
}

// GetAllTasks retrieves all tasks for a user.
func (s *InMemoryStore) GetAllTasks(ownerID int) []Task {
	s.taskMutex.RLock()
	defer s.taskMutex.RUnlock()

	tasks := []Task{}
	for _, task := range s.Tasks {
		if task.OwnerID == ownerID {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

// GenerateToken creates a secure random token.
func GenerateToken() (string, error) {
	b := make([]byte, TokenLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// HashPassword hashes the given password using SHA-256.
// Note: In production, use a stronger hashing algorithm like bcrypt or Argon2.
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return base64.URLEncoding.EncodeToString(hash[:])
}

// CheckPasswordHash compares a plaintext password with a hashed password.
func CheckPasswordHash(password, hash string) bool {
	return HashPassword(password) == hash
}

// Claims represents the JWT claims.
type Claims struct {
	UserID int `json:"user_id"`
}

// GenerateJWT generates a JWT token for a user.
func GenerateJWT(userID int) (string, error) {
	token, err := GenerateToken()
	if err != nil {
		return "", err
	}
	return token, nil
}

// HandlerContext holds the application state.
type HandlerContext struct {
	Store *InMemoryStore
}

// RegisterHandler handles user registration.
func (hc *HandlerContext) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		RespondJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid request payload",
		})
		return
	}

	if req.Username == "" || req.Password == "" {
		RespondJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Username and password are required",
		})
		return
	}

	user, err := hc.Store.AddUser(req.Username, req.Password)
	if err != nil {
		RespondJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// Generate token
	token, err := GenerateJWT(user.ID)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Could not generate token",
		})
		return
	}

	// Store token
	hc.Store.tokenMux.Lock()
	hc.Store.tokenStore[token] = user.ID
	hc.Store.tokenMux.Unlock()

	user.Token = token

	RespondJSON(w, http.StatusCreated, Response{
		Status:  "success",
		Message: "User registered successfully",
		Data: map[string]interface{}{
			"user_id": user.ID,
			"token":   user.Token,
		},
	})
}

// LoginHandler handles user authentication.
func (hc *HandlerContext) LoginHandler(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		RespondJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid request payload",
		})
		return
	}

	user, err := hc.Store.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		RespondJSON(w, http.StatusUnauthorized, Response{
			Status:  "error",
			Message: "Invalid username or password",
		})
		return
	}

	// Generate new token
	token, err := GenerateJWT(user.ID)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Could not generate token",
		})
		return
	}

	// Store token
	hc.Store.tokenMux.Lock()
	hc.Store.tokenStore[token] = user.ID
	hc.Store.tokenMux.Unlock()

	user.Token = token

	RespondJSON(w, http.StatusOK, Response{
		Status:  "success",
		Message: "Logged in successfully",
		Data: map[string]interface{}{
			"user_id": user.ID,
			"token":   user.Token,
		},
	})
}

// CreateTaskHandler handles the creation of a new task.
func (hc *HandlerContext) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		RespondJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid request payload",
		})
		return
	}

	if req.Title == "" {
		RespondJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Title is required",
		})
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value(ContextUserKey).(int)
	if !ok {
		RespondJSON(w, http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Could not retrieve user from context",
		})
		return
	}

	task := hc.Store.AddTask(req.Title, req.Description, userID)

	RespondJSON(w, http.StatusCreated, Response{
		Status:  "success",
		Message: "Task created successfully",
		Data:    task,
	})
}

// GetAllTasksHandler retrieves all tasks for the authenticated user.
func (hc *HandlerContext) GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value(ContextUserKey).(int)
	if !ok {
		RespondJSON(w, http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Could not retrieve user from context",
		})
		return
	}

	tasks := hc.Store.GetAllTasks(userID)

	RespondJSON(w, http.StatusOK, Response{
		Status:  "success",
		Message: "Tasks retrieved successfully",
		Data:    tasks,
	})
}

// GetTaskHandler retrieves a specific task by ID.
func (hc *HandlerContext) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskIDStr, exists := vars["id"]
	if !exists {
		RespondJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Task ID is required",
		})
		return
	}

	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		RespondJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid Task ID",
		})
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value(ContextUserKey).(int)
	if !ok {
		RespondJSON(w, http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Could not retrieve user from context",
		})
		return
	}

	task, err := hc.Store.GetTask(taskID, userID)
	if err != nil {
		RespondJSON(w, http.StatusNotFound, Response{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	RespondJSON(w, http.StatusOK, Response{
		Status:  "success",
		Message: "Task retrieved successfully",
		Data:    task,
	})
}

// UpdateTaskHandler updates an existing task.
func (hc *HandlerContext) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskIDStr, exists := vars["id"]
	if !exists {
		RespondJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Task ID is required",
		})
		return
	}

	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		RespondJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid Task ID",
		})
		return
	}

	type Request struct {
		Title       string `json:"title,omitempty"`
		Description string `json:"description,omitempty"`
		Completed   bool   `json:"completed,omitempty"`
	}

	var req Request
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		RespondJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid request payload",
		})
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value(ContextUserKey).(int)
	if !ok {
		RespondJSON(w, http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Could not retrieve user from context",
		})
		return
	}

	updatedTask := Task{
		Title:       req.Title,
		Description: req.Description,
		Completed:   req.Completed,
	}

	task, err := hc.Store.UpdateTask(taskID, userID, updatedTask)
	if err != nil {
		RespondJSON(w, http.StatusNotFound, Response{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	RespondJSON(w, http.StatusOK, Response{
		Status:  "success",
		Message: "Task updated successfully",
		Data:    task,
	})
}

// DeleteTaskHandler deletes a task by ID.
func (hc *HandlerContext) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskIDStr, exists := vars["id"]
	if !exists {
		RespondJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Task ID is required",
		})
		return
	}

	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		RespondJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid Task ID",
		})
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value(ContextUserKey).(int)
	if !ok {
		RespondJSON(w, http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Could not retrieve user from context",
		})
		return
	}

	err = hc.Store.DeleteTask(taskID, userID)
	if err != nil {
		RespondJSON(w, http.StatusNotFound, Response{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	RespondJSON(w, http.StatusOK, Response{
		Status:  "success",
		Message: "Task deleted successfully",
	})
}

// LoggingMiddleware logs incoming HTTP requests.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		log.Printf("Started %s %s", r.Method, r.RequestURI)

		next.ServeHTTP(w, r)

		duration := time.Since(startTime)
		log.Printf("Completed %s in %v", r.RequestURI, duration)
	})
}

// AuthMiddleware authenticates requests using tokens.
func (hc *HandlerContext) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			RespondJSON(w, http.StatusUnauthorized, Response{
				Status:  "error",
				Message: "Missing Authorization header",
			})
			return
		}

		// Check if the header is in the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			RespondJSON(w, http.StatusUnauthorized, Response{
				Status:  "error",
				Message: "Invalid Authorization header format",
			})
			return
		}

		token := parts[1]

		// Validate token
		hc.Store.tokenMux.RLock()
		userID, exists := hc.Store.tokenStore[token]
		hc.Store.tokenMux.RUnlock()
		if !exists {
			RespondJSON(w, http.StatusUnauthorized, Response{
				Status:  "error",
				Message: "Invalid or expired token",
			})
			return
		}

		// Add user ID to context
		ctx := context.WithValue(r.Context(), ContextUserKey, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// RespondJSON sends a JSON response with the given status code and payload.
func RespondJSON(w http.ResponseWriter, statusCode int, payload Response) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "JSON Marshalling Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

// NotFoundHandler handles undefined routes.
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, http.StatusNotFound, Response{
		Status:  "error",
		Message: "Endpoint not found",
	})
}

// MethodNotAllowedHandler handles invalid HTTP methods.
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, http.StatusMethodNotAllowed, Response{
		Status:  "error",
		Message: "Method not allowed",
	})
}

func main() {
	// Initialize the in-memory store
	store := NewStore()

	// Create a default admin user
	adminUsername := "admin"
	adminPassword := "admin123"
	_, err := store.AddUser(adminUsername, adminPassword)
	if err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	// Initialize the handler context
	hc := &HandlerContext{
		Store: store,
	}

	// Initialize the router
	router := mux.NewRouter().StrictSlash(true)

	// Apply Logging Middleware globally
	router.Use(LoggingMiddleware)

	// Public Routes
	router.HandleFunc("/register", hc.RegisterHandler).Methods("POST")
	router.HandleFunc("/login", hc.LoginHandler).Methods("POST")

	// Protected Routes
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(hc.AuthMiddleware)

	// Task Routes
	protected.HandleFunc("/tasks", hc.CreateTaskHandler).Methods("POST")
	protected.HandleFunc("/tasks", hc.GetAllTasksHandler).Methods("GET")
	protected.HandleFunc("/tasks/{id}", hc.GetTaskHandler).Methods("GET")
	protected.HandleFunc("/tasks/{id}", hc.UpdateTaskHandler).Methods("PUT")
	protected.HandleFunc("/tasks/{id}", hc.DeleteTaskHandler).Methods("DELETE")

	// Handle undefined routes
	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(MethodNotAllowedHandler)

	// Start the server
	port := getEnv("PORT", ServerPort)
	log.Printf("Server is running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// getEnv retrieves the value of the environment variable named by the key.
// It returns the value, or defaultVal if the variable is not present.
func getEnv(key, defaultVal string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		return defaultVal
	}
	return val
}
