// main.go
package cache

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// InMemoryStore holds the users in memory
type InMemoryStore struct {
	sync.RWMutex
	users map[string]User
}

// NewInMemoryStore initializes the in-memory store
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		users: make(map[string]User),
	}
}

// AddUser adds a new user to the store
func (s *InMemoryStore) AddUser(user User) {
	s.Lock()
	defer s.Unlock()
	s.users[user.ID] = user
}

// GetUser retrieves a user by ID
func (s *InMemoryStore) GetUser(id string) (User, bool) {
	s.RLock()
	defer s.RUnlock()
	user, exists := s.users[id]
	return user, exists
}

// UpdateUser updates an existing user
func (s *InMemoryStore) UpdateUser(id string, user User) bool {
	s.Lock()
	defer s.Unlock()
	if _, exists := s.users[id]; exists {
		s.users[id] = user
		return true
	}
	return false
}

// DeleteUser removes a user from the store
func (s *InMemoryStore) DeleteUser(id string) bool {
	s.Lock()
	defer s.Unlock()
	if _, exists := s.users[id]; exists {
		delete(s.users, id)
		return true
	}
	return false
}

// ListUsers lists all users
func (s *InMemoryStore) ListUsers() []User {
	s.RLock()
	defer s.RUnlock()
	users := make([]User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}

// App holds the application configurations
type App struct {
	Router *mux.Router
	Store  *InMemoryStore
}

// Initialize sets up the application
func (a *App) Initialize() {
	a.Store = NewInMemoryStore()
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	a.ApplyMiddleware()
	a.initializeAdditionalRoutes()
}

// initializeRoutes sets up the routes for the API
func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/users", a.createUser).Methods("POST")
	a.Router.HandleFunc("/users", a.getUsers).Methods("GET")
	a.Router.HandleFunc("/users/{id}", a.getUser).Methods("GET")
	a.Router.HandleFunc("/users/{id}", a.updateUser).Methods("PUT")
	a.Router.HandleFunc("/users/{id}", a.deleteUser).Methods("DELETE")
	a.Router.HandleFunc("/", a.home).Methods("GET")
}

// Run starts the HTTP server
func (a *App) Run(addr string) {
	srv := &http.Server{
		Handler:      a.Router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Starting server on %s", addr)
	log.Fatal(srv.ListenAndServe())
}

// home handles the root endpoint
func (a *App) home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to the User API!"))
}

// createUser handles creating a new user
func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Basic validation
	if user.ID == "" || user.FirstName == "" || user.LastName == "" || user.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Missing required user fields")
		return
	}

	// Check if user already exists
	if _, exists := a.Store.GetUser(user.ID); exists {
		respondWithError(w, http.StatusConflict, "User with this ID already exists")
		return
	}

	user.CreatedAt = time.Now()
	a.Store.AddUser(user)
	respondWithJSON(w, http.StatusCreated, user)
}

// getUser handles retrieving a single user by ID
func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, exists := a.Store.GetUser(id)
	if !exists {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

// getUsers handles retrieving all users
func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {
	users := a.Store.ListUsers()
	respondWithJSON(w, http.StatusOK, users)
}

// updateUser handles updating an existing user
func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var user User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Basic validation
	if user.FirstName == "" || user.LastName == "" || user.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Missing required user fields")
		return
	}

	existingUser, exists := a.Store.GetUser(id)
	if !exists {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	// Preserve the original CreatedAt timestamp
	user.ID = id
	user.CreatedAt = existingUser.CreatedAt

	if updated := a.Store.UpdateUser(id, user); !updated {
		respondWithError(w, http.StatusInternalServerError, "Could not update user")
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

// deleteUser handles deleting a user
func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if deleted := a.Store.DeleteUser(id); !deleted {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// respondWithError sends an error response in JSON format
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON sends a response in JSON format
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		// In case of error during marshaling, send a 500 response
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Logger is a middleware for logging HTTP requests
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		log.Printf("Started %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s in %v", r.RequestURI, time.Since(startTime))
	})
}

// Recoverer is a middleware for recovering from panics
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("Recovered from panic: %v", rec)
				respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// ApplyMiddleware applies middleware to the router
func (a *App) ApplyMiddleware() {
	a.Router.Use(Logger)
	a.Router.Use(Recoverer)
}

// initializeAdditionalRoutes sets up additional routes for the API
func (a *App) initializeAdditionalRoutes() {
	a.Router.HandleFunc("/health", a.healthCheck).Methods("GET")
	a.Router.HandleFunc("/version", a.version).Methods("GET")
}

// healthCheck returns the health status of the application
func (a *App) healthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

// version returns the application version
func (a *App) version(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"version": "1.0.0"})
}

func main() {
	app := &App{}
	app.Initialize()
	app.Run(":8080")
}
