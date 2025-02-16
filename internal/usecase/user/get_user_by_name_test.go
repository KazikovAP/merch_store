package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KazikovAP/merch_store/internal/domain"
	userRepo "github.com/KazikovAP/merch_store/internal/repository/user"
	"github.com/KazikovAP/merch_store/internal/usecase/user"
	"github.com/stretchr/testify/assert"
)

// Тест на успешное получение пользователя по имени.
func TestUseCase_GetByUsername_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	username := "alica"
	expectedUser := &domain.User{
		ID:        1,
		Username:  username,
		Balance:   100,
		CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	mock.ExpectQuery(`SELECT id, username, balance, created_at FROM users WHERE username = \$1`).
		WithArgs(username).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "balance", "created_at"}).
			AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Balance, expectedUser.CreatedAt))

	repo := userRepo.NewUserRepository(db)
	useCase := user.NewUseCase(repo)

	u, err := useCase.GetByUsername(context.Background(), username)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, u)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// Тест на ошибку при получении пользователя по имени.
func TestUseCase_GetByUsername_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	username := "bob"

	mock.ExpectQuery(`SELECT id, username, balance, created_at FROM users WHERE username = \$1`).
		WithArgs(username).
		WillReturnError(errors.New("database error"))

	repo := userRepo.NewUserRepository(db)
	useCase := user.NewUseCase(repo)

	u, err := useCase.GetByUsername(context.Background(), username)
	assert.Error(t, err)
	assert.Nil(t, u)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
