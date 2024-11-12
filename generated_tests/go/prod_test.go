package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// MockStore is a mock in-memory store for testing
type MockStore struct {
	books map[int]Book
	users map[int]User
}

func NewMockStore() *MockStore {
	return &MockStore{
		books: make(map[int]Book),
		users: make(map[int]User),
	}
}

func (m *MockStore) AddBook(book Book) Book {
	book.ID = len(m.books) + 1
	book.Available = true
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()
	m.books[book.ID] = book
	return book
}

func (m *MockStore) GetBook(id int) (Book, error) {
	book, exists := m.books[id]
	if !exists {
		return Book{}, errors.New("book not found")
	}
	return book, nil
}

func (m *MockStore) UpdateBook(id int, updated Book) (Book, error) {
	book, exists := m.books[id]
	if !exists {
		return Book{}, errors.New("book not found")
	}
	book.Title = updated.Title
	book.Author = updated.Author
	book.PublishedYear = updated.PublishedYear
	book.ISBN = updated.ISBN
	book.UpdatedAt = time.Now()
	m.books[id] = book
	return book, nil
}

func (m *MockStore) DeleteBook(id int) error {
	_, exists := m.books[id]
	if !exists {
		return errors.New("book not found")
	}
	delete(m.books, id)
	return nil
}

func (m *MockStore) ListBooks() []Book {
	books := []Book{}
	for _, book := range m.books {
		books = append(books, book)
	}
	return books
}

func (m *MockStore) AddUser(user User) User {
	user.ID = len(m.users) + 1
	user.IsActive = true
	user.JoinedAt = time.Now()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	return user
}

func (m *MockStore) GetUser(id int) (User, error) {
	user, exists := m.users[id]
	if !exists {
		return User{}, errors.New("user not found")
	}
	return user, nil
}

func (m *MockStore) UpdateUser(id int, updated User) (User, error) {
	user, exists := m.users[id]
	if !exists {
		return User{}, errors.New("user not found")
	}
	user.Name = updated.Name
	user.Email = updated.Email
	user.IsActive = updated.IsActive
	user.UpdatedAt = time.Now()
	m.users[id] = user
	return user, nil
}

func (m *MockStore) DeleteUser(id int) error {
	_, exists := m.users[id]
	if !exists {
		return errors.New("user not found")
	}
	delete(m.users, id)
	return nil
}

func (m *MockStore) ListUsers() []User {
	users := []User{}
	for _, user := range m.users {
		users = append(users, user)
	}
	return users
}

func TestAddBook(t *testing.T) {
	store := NewMockStore()
	bookService := NewBookService(store)

	book := Book{
		Title:         "Test Book",
		Author:        "Test Author",
		PublishedYear: 2020,
		ISBN:          "123-456-789",
	}

	addedBook := bookService.AddBook(book)

	if addedBook.ID == 0 {
		t.Errorf("Expected book ID to be set, got %d", addedBook.ID)
	}
	if addedBook.Title != book.Title {
		t.Errorf("Expected book title %s, got %s", book.Title, addedBook.Title)
	}
}

func TestGetBookNotFound(t *testing.T) {
	store := NewMockStore()
	bookService := NewBookService(store)

	_, err := bookService.GetBook(1)
	if err == nil {
		t.Errorf("Expected error for non-existent book, got nil")
	}
}

func TestUpdateBookNotFound(t *testing.T) {
	store := NewMockStore()
	bookService := NewBookService(store)

	_, err := bookService.UpdateBook(1, Book{})
	if err == nil {
		t.Errorf("Expected error for non-existent book, got nil")
	}
}

func TestDeleteBookNotFound(t *testing.T) {
	store := NewMockStore()
	bookService := NewBookService(store)

	err := bookService.DeleteBook(1)
	if err == nil {
		t.Errorf("Expected error for non-existent book, got nil")
	}
}

func TestListBooksEmpty(t *testing.T) {
	store := NewMockStore()
	bookService := NewBookService(store)

	books := bookService.ListBooks()
	if len(books) != 0 {
		t.Errorf("Expected no books, got %d", len(books))
	}
}

func TestAddUser(t *testing.T) {
	store := NewMockStore()
	userService := NewUserService(store)

	user := User{
		Name:  "Test User",
		Email: "test@example.com",
	}

	addedUser := userService.AddUser(user)

	if addedUser.ID == 0 {
		t.Errorf("Expected user ID to be set, got %d", addedUser.ID)
	}
	if addedUser.Name != user.Name {
		t.Errorf("Expected user name %s, got %s", user.Name, addedUser.Name)
	}
}

func TestGetUserNotFound(t *testing.T) {
	store := NewMockStore()
	userService := NewUserService(store)

	_, err := userService.GetUser(1)
	if err == nil {
		t.Errorf("Expected error for non-existent user, got nil")
	}
}

func TestUpdateUserNotFound(t *testing.T) {
	store := NewMockStore()
	userService := NewUserService(store)

	_, err := userService.UpdateUser(1, User{})
	if err == nil {
		t.Errorf("Expected error for non-existent user, got nil")
	}
}

func TestDeleteUserNotFound(t *testing.T) {
	store := NewMockStore()
	userService := NewUserService(store)

	err := userService.DeleteUser(1)
	if err == nil {
		t.Errorf("Expected error for non-existent user, got nil")
	}
}

func TestListUsersEmpty(t *testing.T) {
	store := NewMockStore()
	userService := NewUserService(store)

	users := userService.ListUsers()
	if len(users) != 0 {
		t.Errorf("Expected no users, got %d", len(users))
	}
}