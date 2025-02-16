package user_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KazikovAP/merch_store/internal/domain"
	"github.com/KazikovAP/merch_store/internal/repository/user"
	"github.com/stretchr/testify/assert"
)

// Тест на успешное создание пользователя.
func TestRepo_Create_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	username := "alice"
	now := time.Now()
	expectedUser := &domain.User{
		ID:        1,
		Username:  username,
		Balance:   1000,
		CreatedAt: now,
	}

	mock.ExpectQuery(`INSERT INTO users \(username, balance\) VALUES \(\$1, 1000\) RETURNING id, username, balance, created_at`).
		WithArgs(username).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "balance", "created_at"}).
			AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Balance, expectedUser.CreatedAt))

	repo := user.NewUserRepository(db)

	u, err := repo.Create(context.Background(), username)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, u)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// Тест на ошибку при создании пользователя.
func TestRepo_Create_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	username := "alina"

	mock.ExpectQuery(`INSERT INTO users \(username, balance\) VALUES \(\$1, 1000\) RETURNING id, username, balance, created_at`).
		WithArgs(username).
		WillReturnError(errors.New("database error"))

	repo := user.NewUserRepository(db)

	u, err := repo.Create(context.Background(), username)
	assert.Error(t, err)
	assert.Nil(t, u)
	assert.Contains(t, err.Error(), "failed to create user")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// Тест на успешное получение пользователя по ID.
func TestRepo_GetByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	id := 1
	now := time.Now()
	expectedUser := &domain.User{
		ID:        id,
		Username:  "kek",
		Balance:   1000,
		CreatedAt: now,
	}

	mock.ExpectQuery(`SELECT id, username, balance, created_at FROM users WHERE id = \$1`).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "balance", "created_at"}).
			AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Balance, expectedUser.CreatedAt))

	repo := user.NewUserRepository(db)

	u, err := repo.GetByID(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, u)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// Тест на ошибку при получении пользователя по ID (пользователь не найден).
func TestRepo_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	id := 1

	mock.ExpectQuery(`SELECT id, username, balance, created_at FROM users WHERE id = \$1`).
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	repo := user.NewUserRepository(db)

	u, err := repo.GetByID(context.Background(), id)
	assert.Error(t, err)
	assert.Nil(t, u)
	assert.Contains(t, err.Error(), "failed to get user by id")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// Тест на успешное получение пользователя по имени.
func TestRepo_GetByUsername_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	username := "bob"
	now := time.Now()
	expectedUser := &domain.User{
		ID:        1,
		Username:  username,
		Balance:   1000,
		CreatedAt: now,
	}

	mock.ExpectQuery(`SELECT id, username, balance, created_at FROM users WHERE username = \$1`).
		WithArgs(username).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "balance", "created_at"}).
			AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Balance, expectedUser.CreatedAt))

	repo := user.NewUserRepository(db)

	u, err := repo.GetByUsername(context.Background(), username)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, u)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// Тест на ошибку при получении пользователя по имени (пользователь не найден).
func TestRepo_GetByUsername_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	username := "unknown_user"

	mock.ExpectQuery(`SELECT id, username, balance, created_at FROM users WHERE username = \$1`).
		WithArgs(username).
		WillReturnError(sql.ErrNoRows)

	repo := user.NewUserRepository(db)

	u, err := repo.GetByUsername(context.Background(), username)
	assert.Error(t, err)
	assert.Nil(t, u)
	assert.Contains(t, err.Error(), "failed to get user by username")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
