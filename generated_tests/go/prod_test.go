package cache

import (
	"testing"
	"time"
)

func setupMockStore() *MockStore {
	return NewMockStore()
}

func TestMockStore_AddBookSuccess(t *testing.T) {
	store := setupMockStore()
	book := Book{Title: "Effective Go", Author: "Rob Pike", PublishedYear: 2020, ISBN: "987-654-321"}
	addedBook := store.AddBook(book)
	if addedBook.ID == 0 {
		t.Errorf("Book ID not set after adding book")
	}
	if addedBook.Title != book.Title {
		t.Errorf("Expected book title %s, got %s", book.Title, addedBook.Title)
	}
}

func TestMockStore_GetBookSuccess(t *testing.T) {
	store := setupMockStore()
	book := store.AddBook(Book{Title: "Effective Go", Author: "Rob Pike", PublishedYear: 2020, ISBN: "987-654-321"})
	gotBook, err := store.GetBook(book.ID)
	if err != nil {
		t.Errorf("Error getting book: %v", err)
	}
	if gotBook.ID != book.ID {
		t.Errorf("Expected book ID %d, got %d", book.ID, gotBook.ID)
	}
}

func TestMockStore_GetBookFailure(t *testing.T) {
	store := setupMockStore()
	_, err := store.GetBook(1) // Assuming 1 is an ID that doesn't exist
	if err == nil {
		t.Error("Expected error getting non-existent book, got nil")
	}
}

func TestMockStore_UpdateBookSuccess(t *testing.T) {
	store := setupMockStore()
	book := store.AddBook(Book{Title: "Go Programming", Author: "Rob Pike", PublishedYear: 2019, ISBN: "123-456-789"})
	update := Book{Title: "Advanced Go Programming", Author: "Rob Pike", PublishedYear: 2020, ISBN: "123-456-789"}
	updatedBook, err := store.UpdateBook(book.ID, update)
	if err != nil {
		t.Errorf("Error updating book: %v", err)
	}
	if updatedBook.Title != update.Title {
		t.Errorf("Expected updated book title %s, got %s", update.Title, updatedBook.Title)
	}
}

func TestMockStore_UpdateBookFailure(t *testing.T) {
	store := setupMockStore()
	_, err := store.UpdateBook(1, Book{}) // Assuming 1 is an ID that doesn't exist
	if err == nil {
		t.Error("Expected error updating non-existent book, got nil")
	}
}

func TestMockStore_DeleteBookSuccess(t *testing.T) {
	store := setupMockStore()
	book := store.AddBook(Book{Title: "Concurrent Programming in Go", Author: "Katherine Cox-Buday", PublishedYear: 2021, ISBN: "123-321-123"})
	err := store.DeleteBook(book.ID)
	if err != nil {
		t.Errorf("Error deleting book: %v", err)
	}
}

func TestMockStore_DeleteBookFailure(t *testing.T) {
	store := setupMockStore()
	err := store.DeleteBook(1) // Assuming 1 is an ID that doesn't exist
	if err == nil {
		t.Error("Expected error deleting non-existent book, got nil")
	}
}

func TestMockStore_ListBooks(t *testing.T) {
	store := setupMockStore()
	store.AddBook(Book{Title: "Go in Action", Author: "William Kennedy", PublishedYear: 2015, ISBN: "123-456-789"})
	store.AddBook(Book{Title: "Go Programming Language", Author: "Alan A. A. Donovan", PublishedYear: 2016, ISBN: "987-654-321"})
	books := store.ListBooks()
	if len(books) != 2 {
		t.Errorf("Expected 2 books, got %d", len(books))
	}
}

func TestMockStore_AddUserSuccess(t *testing.T) {
	store := setupMockStore()
	user := User{Name: "John Doe", Email: "john@example.com"}
	addedUser := store.AddUser(user)
	if addedUser.ID == 0 {
		t.Errorf("User ID not set after adding user")
	}
	if addedUser.Name != user.Name {
		t.Errorf("Expected user name %s, got %s", user.Name, addedUser.Name)
	}
}

func TestMockStore_GetUserSuccess(t *testing.T) {
	store := setupMockStore()
	user := store.AddUser(User{Name: "Jane Doe", Email: "jane@example.com"})
	gotUser, err := store.GetUser(user.ID)
	if err != nil {
		t.Errorf("Error getting user: %v", err)
	}
	if gotUser.ID != user.ID {
		t.Errorf("Expected user ID %d, got %d", user.ID, gotUser.ID)
	}
}

func TestMockStore_GetUserFailure(t *testing.T) {
	store := setupMockStore()
	_, err := store.GetUser(1) // Assuming 1 is an ID that doesn't exist
	if err == nil {
		t.Error("Expected error getting non-existent user, got nil")
	}
}

func TestMockStore_UpdateUserSuccess(t *testing.T) {
	store := setupMockStore()
	user := store.AddUser(User{Name: "John Doe", Email: "john@example.com"})
	update := User{Name: "Johnathan Doe", Email: "johnathan@example.com"}
	updatedUser, err := store.UpdateUser(user.ID, update)
	if err != nil {
		t.Errorf("Error updating user: %v", err)
	}
	if updatedUser.Name != update.Name {
		t.Errorf("Expected updated user name %s, got %s", update.Name, updatedUser.Name)
	}
}

func TestMockStore_UpdateUserFailure(t *testing.T) {
	store := setupMockStore()
	_, err := store.UpdateUser(1, User{}) // Assuming 1 is an ID that doesn't exist
	if err == nil {
		t.Error("Expected error updating non-existent user, got nil")
	}
}

func TestMockStore_DeleteUserSuccess(t *testing.T) {
	store := setupMockStore()
	user := store.AddUser(User{Name: "Delete Me", Email: "delete@example.com"})
	err := store.DeleteUser(user.ID)
	if err != nil {
		t.Errorf("Error deleting user: %v", err)
	}
}

func TestMockStore_DeleteUserFailure(t *testing.T) {
	store := setupMockStore()
	err := store.DeleteUser(1) // Assuming 1 is an ID that doesn't exist
	if err == nil {
		t.Error("Expected error deleting non-existent user, got nil")
	}
}

func TestMockStore_ListUsers(t *testing.T) {
	store := setupMockStore()
	store.AddUser(User{Name: "User One", Email: "one@example.com"})
	store.AddUser(User{Name: "User Two", Email: "two@example.com"})
	users := store.ListUsers()
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
}