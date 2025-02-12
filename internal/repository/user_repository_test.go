package repository_test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KazikovAP/merch_store/internal/model"
	"github.com/KazikovAP/merch_store/internal/repository"
)

func TestGetByUsername_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "username", "password", "coins"}).
		AddRow(1, "testuser", "hashedpassword", 1000)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, username, password, coins FROM users WHERE username=$1")).
		WithArgs("testuser").
		WillReturnRows(rows)

	userRepo := repository.NewUserRepository(db)

	user, err := userRepo.GetByUsername("testuser")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if user.Username != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", user.Username)
	}
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO users (username, password, coins) VALUES ($1, $2, $3) RETURNING id")).
		WithArgs("testuser", "hashedpassword", 1000).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	userRepo := repository.NewUserRepository(db)

	user := &model.User{
		Username: "testuser",
		Password: "hashedpassword",
		Coins:    1000,
	}

	err = userRepo.Create(user)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if user.ID != 1 {
		t.Errorf("expected user ID 1, got %d", user.ID)
	}
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET password=$1, coins=$2 WHERE id=$3")).
		WithArgs("newhashedpassword", 800, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	userRepo := repository.NewUserRepository(db)

	user := &model.User{
		ID:       1,
		Username: "testuser",
		Password: "newhashedpassword",
		Coins:    800,
	}

	err = userRepo.Update(user)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}
