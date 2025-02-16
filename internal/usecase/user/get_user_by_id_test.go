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

// Тест на успешное получение пользователя по ID.
func TestUseCase_GetByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	userID := 1
	expectedUser := &domain.User{
		ID:        userID,
		Username:  "alice",
		Balance:   100,
		CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	mock.ExpectQuery(`SELECT id, username, balance, created_at FROM users WHERE id = \$1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "balance", "created_at"}).
			AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Balance, expectedUser.CreatedAt))

	repo := userRepo.NewUserRepository(db)
	useCase := user.NewUseCase(repo)

	u, err := useCase.GetByID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, u)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// Тест на ошибку при получении пользователя по ID.
func TestUseCase_GetByID_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	userID := 1

	mock.ExpectQuery(`SELECT id, username, balance, created_at FROM users WHERE id = \$1`).
		WithArgs(userID).
		WillReturnError(errors.New("database error"))

	repo := userRepo.NewUserRepository(db)
	useCase := user.NewUseCase(repo)

	u, err := useCase.GetByID(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, u)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
