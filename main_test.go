package main

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

type User struct {
	ID   int
	Name string
}

type UserRepository interface {
	GetUserByID(id int) (*User, error)
	CreateUser(user *User) error
}

type DBUserRepository struct {
	DB *sql.DB
}

func (repo *DBUserRepository) GetUserByID(id int) (*User, error) {
	var user User
	row := repo.DB.QueryRow("SELECT id, name FROM users WHERE id = ?", id)
	err := row.Scan(&user.ID, &user.Name)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *DBUserRepository) CreateUser(user *User) error {
	_, err := repo.DB.Exec("INSERT INTO users (id, name) VALUES (?, ?)", user.ID, user.Name)
	return err
}

// Positive case
func TestDBUserRepository_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := &DBUserRepository{DB: db}

	mock.ExpectExec("INSERT INTO users").WithArgs(1, "John Doe").WillReturnResult(sqlmock.NewResult(1, 1))

	user := &User{ID: 1, Name: "John Doe"}
	err = repo.CreateUser(user)
	if err != nil {
		t.Errorf("error was not expected while creating user: %s", err)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// Negative Case
func TestDBUserRepository_CreateUserNegative(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := &DBUserRepository{DB: db}

	mock.ExpectExec("INSERT INTO users").WithArgs(2, "John Doe").WillReturnResult(sqlmock.NewResult(1, 1))

	// Use different name
	user := &User{ID: 2, Name: "Jane Doe"}
	err = repo.CreateUser(user)
	if err != nil {
		t.Errorf("error was not expected while creating user: %s", err)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDBUserRepository_GetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := &DBUserRepository{DB: db}

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "John Doe")
	mock.ExpectQuery("SELECT id, name FROM users WHERE id = ?").WithArgs(1).WillReturnRows(rows)

	user, err := repo.GetUserByID(1)
	if err != nil {
		t.Errorf("error was not expected while getting user by ID: %s", err)
	}

	if user.Name != "John Doe" {
		t.Errorf("expected name to be 'John Doe', but got '%s'", user.Name)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
