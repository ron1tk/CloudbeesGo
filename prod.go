package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ---------------------- Models ----------------------

// Book represents a book in the library
type Book struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	PublishedYear int       `json:"published_year"`
	ISBN          string    `json:"isbn"`
	Available     bool      `json:"available"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// User represents a library user
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	JoinedAt  time.Time `json:"joined_at"`
	IsActive  bool      `json:"is_active"`
	Borrowed  []int     `json:"borrowed_books"` // Slice of Book IDs
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ---------------------- In-Memory Store ----------------------

type Store struct {
	books      map[int]Book
	users      map[int]User
	bookMutex  sync.RWMutex
	userMutex  sync.RWMutex
	nextBookID int
	nextUserID int
}

// NewStore initializes the in-memory store
func NewStore() *Store {
	return &Store{
		books:      make(map[int]Book),
		users:      make(map[int]User),
		nextBookID: 1,
		nextUserID: 1,
	}
}

// ---------------------- Services ----------------------

// BookService handles book-related operations
type BookService struct {
	store *Store
}

// NewBookService creates a new BookService
func NewBookService(store *Store) *BookService {
	return &BookService{store: store}
}

// AddBook adds a new book to the store
func (bs *BookService) AddBook(book Book) Book {
	bs.store.bookMutex.Lock()
	defer bs.store.bookMutex.Unlock()
	book.ID = bs.store.nextBookID
	bs.store.nextBookID++
	book.Available = true
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()
	bs.store.books[book.ID] = book
	return book
}

// GetBook retrieves a book by ID
func (bs *BookService) GetBook(id int) (Book, error) {
	bs.store.bookMutex.RLock()
	defer bs.store.bookMutex.RUnlock()
	book, exists := bs.store.books[id]
	if !exists {
		return Book{}, errors.New("book not found")
	}
	return book, nil
}

// UpdateBook updates an existing book
func (bs *BookService) UpdateBook(id int, updated Book) (Book, error) {
	bs.store.bookMutex.Lock()
	defer bs.store.bookMutex.Unlock()
	book, exists := bs.store.books[id]
	if !exists {
		return Book{}, errors.New("book not found")
	}
	if updated.Title != "" {
		book.Title = updated.Title
	}
	if updated.Author != "" {
		book.Author = updated.Author
	}
	if updated.PublishedYear != 0 {
		book.PublishedYear = updated.PublishedYear
	}
	if updated.ISBN != "" {
		book.ISBN = updated.ISBN
	}
	book.UpdatedAt = time.Now()
	bs.store.books[id] = book
	return book, nil
}

// DeleteBook removes a book from the store
func (bs *BookService) DeleteBook(id int) error {
	bs.store.bookMutex.Lock()
	defer bs.store.bookMutex.Unlock()
	_, exists := bs.store.books[id]
	if !exists {
		return errors.New("book not found")
	}
	delete(bs.store.books, id)
	return nil
}

// ListBooks returns all books
func (bs *BookService) ListBooks() []Book {
	bs.store.bookMutex.RLock()
	defer bs.store.bookMutex.RUnlock()
	books := []Book{}
	for _, book := range bs.store.books {
		books = append(books, book)
	}
	return books
}

// UserService handles user-related operations
type UserService struct {
	store *Store
}

// NewUserService creates a new UserService
func NewUserService(store *Store) *UserService {
	return &UserService{store: store}
}

// AddUser adds a new user to the store
func (us *UserService) AddUser(user User) User {
	us.store.userMutex.Lock()
	defer us.store.userMutex.Unlock()
	user.ID = us.store.nextUserID
	us.store.nextUserID++
	user.IsActive = true
	user.JoinedAt = time.Now()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	us.store.users[user.ID] = user
	return user
}

// GetUser retrieves a user by ID
func (us *UserService) GetUser(id int) (User, error) {
	us.store.userMutex.RLock()
	defer us.store.userMutex.RUnlock()
	user, exists := us.store.users[id]
	if !exists {
		return User{}, errors.New("user not found")
	}
	return user, nil
}

// UpdateUser updates an existing user
func (us *UserService) UpdateUser(id int, updated User) (User, error) {
	us.store.userMutex.Lock()
	defer us.store.userMutex.Unlock()
	user, exists := us.store.users[id]
	if !exists {
		return User{}, errors.New("user not found")
	}
	if updated.Name != "" {
		user.Name = updated.Name
	}
	if updated.Email != "" {
		user.Email = updated.Email
	}
	user.IsActive = updated.IsActive
	user.UpdatedAt = time.Now()
	us.store.users[id] = user
	return user, nil
}

// DeleteUser removes a user from the store
func (us *UserService) DeleteUser(id int) error {
	us.store.userMutex.Lock()
	defer us.store.userMutex.Unlock()
	_, exists := us.store.users[id]
	if !exists {
		return errors.New("user not found")
	}
	delete(us.store.users, id)
	return nil
}

// ListUsers returns all users
func (us *UserService) ListUsers() []User {
	us.store.userMutex.RLock()
	defer us.store.userMutex.RUnlock()
	users := []User{}
	for _, user := range us.store.users {
		users = append(users, user)
	}
	return users
}

// BorrowBook allows a user to borrow a book
func (us *UserService) BorrowBook(userID, bookID int, bs *BookService) error {
	us.store.userMutex.Lock()
	defer us.store.userMutex.Unlock()

	book, err := bs.GetBook(bookID)
	if err != nil {
		return err
	}

	if !book.Available {
		return errors.New("book is not available")
	}

	user, exists := us.store.users[userID]
	if !exists {
		return errors.New("user not found")
	}

	// Check if user already borrowed the book
	for _, bID := range user.Borrowed {
		if bID == bookID {
			return errors.New("user already borrowed this book")
		}
	}

	// Update book availability
	book.Available = false
	book.UpdatedAt = time.Now()
	bs.store.bookMutex.Lock()
	bs.store.books[bookID] = book
	bs.store.bookMutex.Unlock()

	// Update user's borrowed books
	user.Borrowed = append(user.Borrowed, bookID)
	user.UpdatedAt = time.Now()
	us.store.users[userID] = user

	return nil
}

// ReturnBook allows a user to return a borrowed book
func (us *UserService) ReturnBook(userID, bookID int, bs *BookService) error {
	us.store.userMutex.Lock()
	defer us.store.userMutex.Unlock()

	user, exists := us.store.users[userID]
	if !exists {
		return errors.New("user not found")
	}

	// Check if user has borrowed the book
	found := false
	for i, bID := range user.Borrowed {
		if bID == bookID {
			// Remove the book from borrowed list
			user.Borrowed = append(user.Borrowed[:i], user.Borrowed[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return errors.New("user has not borrowed this book")
	}

	// Update user
	user.UpdatedAt = time.Now()
	us.store.users[userID] = user

	// Update book availability
	book, err := bs.GetBook(bookID)
	if err != nil {
		return err
	}
	book.Available = true
	book.UpdatedAt = time.Now()
	bs.store.bookMutex.Lock()
	bs.store.books[bookID] = book
	bs.store.bookMutex.Unlock()

	return nil
}

// ---------------------- Handlers ----------------------

// Handler struct contains services
type Handler struct {
	bookService *BookService
	userService *UserService
}

// NewHandler creates a new Handler
func NewHandler(bs *BookService, us *UserService) *Handler {
	return &Handler{
		bookService: bs,
		userService: us,
	}
}

// ServeHTTP implements http.Handler
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Simple routing based on URL path and method
	path := r.URL.Path
	method := r.Method

	// Books endpoints
	if strings.HasPrefix(path, "/books") {
		h.handleBooks(w, r)
		return
	}

	// Users endpoints
	if strings.HasPrefix(path, "/users") {
		h.handleUsers(w, r)
		return
	}

	// Borrow book
	if strings.HasPrefix(path, "/borrow") && method == http.MethodPost {
		h.handleBorrow(w, r)
		return
	}

	// Return book
	if strings.HasPrefix(path, "/return") && method == http.MethodPost {
		h.handleReturn(w, r)
		return
	}

	// Not found
	http.NotFound(w, r)
}

// ---------------------- Book Handlers ----------------------

// handleBooks handles all /books endpoints
func (h *Handler) handleBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listBooks(w, r)
	case http.MethodPost:
		h.createBook(w, r)
	case http.MethodPut:
		h.updateBook(w, r)
	case http.MethodDelete:
		h.deleteBook(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// listBooks handles GET /books
func (h *Handler) listBooks(w http.ResponseWriter, r *http.Request) {
	books := h.bookService.ListBooks()
	respondJSON(w, http.StatusOK, books)
}

// createBook handles POST /books
func (h *Handler) createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if book.Title == "" || book.Author == "" || book.ISBN == "" || book.PublishedYear == 0 {
		respondError(w, http.StatusBadRequest, "Missing required fields")
		return
	}
	created := h.bookService.AddBook(book)
	respondJSON(w, http.StatusCreated, created)
}

// updateBook handles PUT /books?id={id}
func (h *Handler) updateBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "Missing book ID")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}
	var book Book
	err = json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	updated, err := h.bookService.UpdateBook(id, book)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, updated)
}

// deleteBook handles DELETE /books?id={id}
func (h *Handler) deleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "Missing book ID")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}
	err = h.bookService.DeleteBook(id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "Book deleted"})
}

// ---------------------- User Handlers ----------------------

// handleUsers handles all /users endpoints
func (h *Handler) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listUsers(w, r)
	case http.MethodPost:
		h.createUser(w, r)
	case http.MethodPut:
		h.updateUser(w, r)
	case http.MethodDelete:
		h.deleteUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// listUsers handles GET /users
func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) {
	users := h.userService.ListUsers()
	respondJSON(w, http.StatusOK, users)
}

// createUser handles POST /users
func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if user.Name == "" || user.Email == "" {
		respondError(w, http.StatusBadRequest, "Missing required fields")
		return
	}
	created := h.userService.AddUser(user)
	respondJSON(w, http.StatusCreated, created)
}

// updateUser handles PUT /users?id={id}
func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "Missing user ID")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	var user User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	updated, err := h.userService.UpdateUser(id, user)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, updated)
}

// deleteUser handles DELETE /users?id={id}
func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "Missing user ID")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	err = h.userService.DeleteUser(id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "User deleted"})
}

// ---------------------- Borrow & Return Handlers ----------------------

// handleBorrow handles POST /borrow
func (h *Handler) handleBorrow(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID int `json:"user_id"`
		BookID int `json:"book_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if req.UserID == 0 || req.BookID == 0 {
		respondError(w, http.StatusBadRequest, "Missing user_id or book_id")
		return
	}
	err = h.userService.BorrowBook(req.UserID, req.BookID, h.bookService)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "Book borrowed successfully"})
}

// handleReturn handles POST /return
func (h *Handler) handleReturn(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID int `json:"user_id"`
		BookID int `json:"book_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if req.UserID == 0 || req.BookID == 0 {
		respondError(w, http.StatusBadRequest, "Missing user_id or book_id")
		return
	}
	err = h.userService.ReturnBook(req.UserID, req.BookID, h.bookService)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "Book returned successfully"})
}

// ---------------------- Utility Functions ----------------------

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// respondError sends an error response in JSON
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// ---------------------- Sample Data Initialization ----------------------

// initializeSampleData adds some sample books and users to the store
func initializeSampleData(bs *BookService, us *UserService) {
	// Sample Books
	bs.AddBook(Book{
		Title:         "The Great Gatsby",
		Author:        "F. Scott Fitzgerald",
		PublishedYear: 1925,
		ISBN:          "9780743273565",
	})
	bs.AddBook(Book{
		Title:         "1984",
		Author:        "George Orwell",
		PublishedYear: 1949,
		ISBN:          "9780451524935",
	})
	bs.AddBook(Book{
		Title:         "To Kill a Mockingbird",
		Author:        "Harper Lee",
		PublishedYear: 1960,
		ISBN:          "9780061120084",
	})

	// Sample Users
	us.AddUser(User{
		Name:  "Alice Johnson",
		Email: "alice@example.com",
	})
	us.AddUser(User{
		Name:  "Bob Smith",
		Email: "bob@example.com",
	})
	us.AddUser(User{
		Name:  "Charlie Brown",
		Email: "charlie@example.com",
	})
}

// ---------------------- Main Function ----------------------

func main() {
	store := NewStore()
	bookService := NewBookService(store)
	userService := NewUserService(store)
	handler := NewHandler(bookService, userService)

	// Initialize sample data
	initializeSampleData(bookService, userService)

	// Start HTTP server
	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

// ---------------------- Additional Functions ----------------------

// You can add more functions, methods, and types below to expand the functionality.
// For example, adding search functionality, pagination, authentication, etc.

// SearchBooks allows searching books by title or author
func (h *Handler) searchBooks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		respondError(w, http.StatusBadRequest, "Missing search query")
		return
	}
	books := h.bookService.ListBooks()
	matched := []Book{}
	for _, book := range books {
		if strings.Contains(strings.ToLower(book.Title), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(book.Author), strings.ToLower(query)) {
			matched = append(matched, book)
		}
	}
	respondJSON(w, http.StatusOK, matched)
}

// Implement other features as needed to expand the codebase
