// main.go
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/ulule/limiter/v3"
	middlewareLimiter "github.com/ulule/limiter/v3/drivers/middleware/gorilla"
	rateLimitMemory "github.com/ulule/limiter/v3/drivers/store/memory"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/go-playground/validator/v10"
	"github.com/rs/cors"
	"github.com/kelseyhightower/envconfig"
)

// Config holds the application configuration
type Config struct {
	ServerPort      string `envconfig:"SERVER_PORT" default:"8080"`
	JWTSecret       string `envconfig:"JWT_SECRET" default:"mysecretkey"`
	RateLimit       string `envconfig:"RATE_LIMIT" default:"100-M"`
	DatabasePath    string `envconfig:"DATABASE_PATH" default:"users.db"`
	LogLevel        string `envconfig:"LOG_LEVEL" default:"info"`
	TokenExpiryMins int    `envconfig:"TOKEN_EXPIRY_MINS" default:"60"`
}

var cfg Config

// Initialize configuration from environment variables
func initConfig() error {
	return envconfig.Process("", &cfg)
}

// User represents a user in the system
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	FirstName string    `json:"first_name" validate:"required"`
	LastName  string    `json:"last_name" validate:"required"`
	Email     string    `json:"email" gorm:"unique" validate:"required,email"`
	CreatedAt time.Time `json:"created_at"`
}

// UserInput represents the input data for creating/updating a user
type UserInput struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
}

// Credentials represents the login credentials
type Credentials struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Claims represents the JWT claims
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// App holds the application configurations and dependencies
type App struct {
	Router      *mux.Router
	DB          *gorm.DB
	Logger      *logrus.Logger
	Validator   *validator.Validate
	JWTSecret   string
	TokenExpiry time.Duration
}

// Initialize sets up the application
func (a *App) Initialize() error {
	// Initialize Logger
	a.Logger = logrus.New()
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return err
	}
	a.Logger.SetLevel(level)
	a.Logger.SetFormatter(&logrus.JSONFormatter{})

	// Initialize Database
	a.DB, err = gorm.Open(sqlite.Open(cfg.DatabasePath), &gorm.Config{})
	if err != nil {
		return err
	}

	// Migrate the schema
	err = a.DB.AutoMigrate(&User{})
	if err != nil {
		return err
	}

	// Initialize Validator
	a.Validator = validator.New()

	// Initialize Router
	a.Router = mux.NewRouter()

	// Initialize JWT settings
	a.JWTSecret = cfg.JWTSecret
	a.TokenExpiry = time.Duration(cfg.TokenExpiryMins) * time.Minute

	// Initialize Routes
	a.initializeRoutes()

	return nil
}

// initializeRoutes sets up the routes for the API
func (a *App) initializeRoutes() {
	// Public Routes
	a.Router.HandleFunc("/login", a.login).Methods("POST")
	a.Router.HandleFunc("/health", a.healthCheck).Methods("GET")
	a.Router.HandleFunc("/version", a.version).Methods("GET")

	// Protected Routes
	api := a.Router.PathPrefix("/api").Subrouter()
	api.Use(a.jwtMiddleware)
	api.HandleFunc("/users", a.createUser).Methods("POST")
	api.HandleFunc("/users", a.getUsers).Methods("GET")
	api.HandleFunc("/users/{id:[0-9]+}", a.getUser).Methods("GET")
	api.HandleFunc("/users/{id:[0-9]+}", a.updateUser).Methods("PUT")
	api.HandleFunc("/users/{id:[0-9]+}", a.deleteUser).Methods("DELETE")
}

// Run starts the HTTP server with middleware and graceful shutdown
func (a *App) Run(addr string) {
	// Setup Rate Limiter
	rate, err := limiter.NewRateFromFormatted(cfg.RateLimit)
	if err != nil {
		a.Logger.Fatalf("Invalid rate limit format: %v", err)
	}
	store := rateLimitMemory.NewStore()
	instance := limiter.New(store, rate)
	a.Router.Use(middlewareLimiter.NewMiddleware(instance))

	// Setup CORS
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Adjust as needed
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})
	handler := corsMiddleware.Handler(a.Router)

	// Create HTTP Server
	srv := &http.Server{
		Handler:      handler,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Start Server in a Goroutine
	go func() {
		a.Logger.Infof("Starting server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.Logger.Fatalf("Could not listen on %s: %v\n", addr, err)
		}
	}()

	// Graceful Shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	a.Logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	srv.Shutdown(ctx)

	a.Logger.Info("Server gracefully stopped")
}

// jwtMiddleware protects routes using JWT authentication
func (a *App) jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			respondWithError(w, http.StatusUnauthorized, "Invalid Authorization header format")
			return
		}

		tokenStr := parts[1]

		// Parse the token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(a.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Token is valid, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

// login handles user authentication and token generation
func (a *App) login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate input
	err = a.Validator.Struct(creds)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Missing or invalid fields")
		return
	}

	// For demonstration, use hardcoded credentials
	if creds.Username != "admin" || creds.Password != "password" {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Create JWT Token
	expirationTime := time.Now().Add(a.TokenExpiry)
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "user-management-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(a.JWTSecret))
	if err != nil {
		a.Logger.Errorf("Error signing token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Could not generate token")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"token": tokenStr})
}

// createUser handles creating a new user
func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	var input UserInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate input
	err = a.Validator.Struct(input)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Missing or invalid fields")
		return
	}

	user := User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		CreatedAt: time.Now(),
	}

	result := a.DB.Create(&user)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "UNIQUE constraint failed") {
			respondWithError(w, http.StatusConflict, "Email already exists")
			return
		}
		a.Logger.Errorf("Error creating user: %v", result.Error)
		respondWithError(w, http.StatusInternalServerError, "Could not create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

// getUser handles retrieving a single user by ID
func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, _ := strconv.Atoi(idStr)

	var user User
	result := a.DB.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	} else if result.Error != nil {
		a.Logger.Errorf("Error retrieving user: %v", result.Error)
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve user")
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

// getUsers handles retrieving all users with pagination and filtering
func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {
	// Pagination parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Filtering parameters
	firstName := r.URL.Query().Get("first_name")
	lastName := r.URL.Query().Get("last_name")
	email := r.URL.Query().Get("email")

	var users []User
	query := a.DB.Model(&User{})

	if firstName != "" {
		query = query.Where("first_name LIKE ?", "%"+firstName+"%")
	}
	if lastName != "" {
		query = query.Where("last_name LIKE ?", "%"+lastName+"%")
	}
	if email != "" {
		query = query.Where("email LIKE ?", "%"+email+"%")
	}

	result := query.Limit(limit).Offset(offset).Find(&users)
	if result.Error != nil {
		a.Logger.Errorf("Error retrieving users: %v", result.Error)
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve users")
		return
	}

	respondWithJSON(w, http.StatusOK, users)
}

// updateUser handles updating an existing user
func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, _ := strconv.Atoi(idStr)

	var input UserInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate input
	err = a.Validator.Struct(input)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Missing or invalid fields")
		return
	}

	var user User
	result := a.DB.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	} else if result.Error != nil {
		a.Logger.Errorf("Error retrieving user: %v", result.Error)
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve user")
		return
	}

	user.FirstName = input.FirstName
	user.LastName = input.LastName
	user.Email = input.Email

	result = a.DB.Save(&user)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "UNIQUE constraint failed") {
			respondWithError(w, http.StatusConflict, "Email already exists")
			return
		}
		a.Logger.Errorf("Error updating user: %v", result.Error)
		respondWithError(w, http.StatusInternalServerError, "Could not update user")
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

// deleteUser handles deleting a user
func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, _ := strconv.Atoi(idStr)

	result := a.DB.Delete(&User{}, id)
	if result.Error != nil {
		a.Logger.Errorf("Error deleting user: %v", result.Error)
		respondWithError(w, http.StatusInternalServerError, "Could not delete user")
		return
	}
	if result.RowsAffected == 0 {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// healthCheck returns the health status of the application
func (a *App) healthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

// version returns the application version
func (a *App) version(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"version": "1.0.0"})
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

func main() {
	// Initialize Configuration
	err := initConfig()
	if err != nil {
		fmt.Printf("Error initializing config: %v\n", err)
		os.Exit(1)
	}

	// Initialize App
	app := &App{}
	err = app.Initialize()
	if err != nil {
		app.Logger.Fatalf("Failed to initialize the application: %v", err)
	}

	// Run the App
	app.Run(":" + cfg.ServerPort)
}
