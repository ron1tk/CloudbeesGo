package main

import (
	"testing"
	"time"
)

func setupMockStoreWithBook() (*MockStore, Book) {
	mockStore := NewMockStore()
	testBook := Book{
		Title:         "Existing Book",
		Author:        "Existing Author",
		PublishedYear: 1999,
		ISBN:          "999-999-999",
	}
	mockStore.AddBook(testBook)
	return mockStore, testBook
}

func setupMockStoreWithUser() (*MockStore, User) {
	mockStore := NewMockStore()
	testUser := User{
		Name:  "Existing User",
		Email: "existing@example.com",
	}
	mockStore.AddUser(testUser)
	return mockStore, testUser
}

func TestAddBookSuccess(t *testing.T) {
	store := NewMockStore()
	bookService := NewBookService(store)

	book := Book{
		Title:         "New Book",
		Author:        "New Author",
		PublishedYear: 2021,
		ISBN:          "321-654-987",
	}

	addedBook := bookService.AddBook(book)

	if addedBook.ID != 1 {
		t.Errorf("Expected book ID to be 1, got %d", addedBook.ID)
	}
	if addedBook.Title != book.Title {
		t.Errorf("Expected book title %s, got %s", book.Title, addedBook.Title)
	}
	if !addedBook.Available {
		t.Errorf("Expected book to be available, but it was not")
	}
}

func TestGetBookSuccess(t *testing.T) {
	store, testBook := setupMockStoreWithBook()
	bookService := NewBookService(store)

	foundBook, err := bookService.GetBook(testBook.ID)
	if err != nil {
		t.Errorf("Did not expect an error, but got %s", err)
	}
	if foundBook.ID != testBook.ID {
		t.Errorf("Expected book ID %d, got %d", testBook.ID, foundBook.ID)
	}
}

func TestUpdateBookSuccess(t *testing.T) {
	store, testBook := setupMockStoreWithBook()
	bookService := NewBookService(store)

	update := Book{
		Title: "Updated Book",
		Author: "Updated Author",
		PublishedYear: 2001,
		ISBN: "111-222-333",
	}

	updatedBook, err := bookService.UpdateBook(testBook.ID, update)
	if err != nil {
		t.Errorf("Did not expect an error, but got %s", err)
	}
	if updatedBook.Title != update.Title {
		t.Errorf("Expected book title %s, got %s", update.Title, updatedBook.Title)
	}
}

func TestDeleteBookSuccess(t *testing.T) {
	store, testBook := setupMockStoreWithBook()
	bookService := NewBookService(store)

	err := bookService.DeleteBook(testBook.ID)
	if err != nil {
		t.Errorf("Did not expect an error, but got %s", err)
	}
	if _, err := bookService.GetBook(testBook.ID); err == nil {
		t.Errorf("Expected an error for non-existent book, got nil")
	}
}

func TestListBooksNotEmpty(t *testing.T) {
	store, _ := setupMockStoreWithBook()
	bookService := NewBookService(store)

	books := bookService.ListBooks()
	if len(books) == 0 {
		t.Errorf("Expected books, got none")
	}
}

func TestAddUserSuccess(t *testing.T) {
	store := NewMockStore()
	userService := NewUserService(store)

	user := User{
		Name:  "New User",
		Email: "new@example.com",
	}

	addedUser := userService.AddUser(user)

	if addedUser.ID == 0 {
		t.Errorf("Expected user ID to be set, got %d", addedUser.ID)
	}
	if addedUser.Name != user.Name {
		t.Errorf("Expected user name %s, got %s", user.Name, addedUser.Name)
	}
	if !addedUser.IsActive {
		t.Errorf("Expected user to be active, but it was not")
	}
}

func TestGetUserSuccess(t *testing.T) {
	store, testUser := setupMockStoreWithUser()
	userService := NewUserService(store)

	foundUser, err := userService.GetUser(testUser.ID)
	if err != nil {
		t.Fatalf("Did not expect an error, but got %s", err)
	}
	if foundUser.ID != testUser.ID {
		t.Errorf("Expected user ID %d, got %d", testUser.ID, foundUser.ID)
	}
}

func TestUpdateUserSuccess(t *testing.T) {
	store, testUser := setupMockStoreWithUser()
	userService := NewUserService(store)

	update := User{
		Name: "Updated User",
		Email: "updated@example.com",
		IsActive: false,
	}

	updatedUser, err := userService.UpdateUser(testUser.ID, update)
	if err != nil {
		t.Errorf("Did not expect an error, but got %s", err)
	}
	if updatedUser.Name != update.Name {
		t.Errorf("Expected user name %s, got %s", update.Name, updatedUser.Name)
	}
}

func TestDeleteUserSuccess(t *testing.T) {
	store, testUser := setupMockStoreWithUser()
	userService := NewUserService(store)

	err := userService.DeleteUser(testUser.ID)
	if err != nil {
		t.Errorf("Did not expect an error, but got %s", err)
	}
	if _, err := userService.GetUser(testUser.ID); err == nil {
		t.Errorf("Expected an error for non-existent user, got nil")
	}
}

func TestListUsersNotEmpty(t *testing.T) {
	store, _ := setupMockStoreWithUser()
	userService := NewUserService(store)

	users := userService.ListUsers()
	if len(users) == 0 {
		t.Errorf("Expected users, got none")
	}
}