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

// Тест на успешную регистрацию пользователя.
func TestUseCase_Register_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	username := "alice"
	expectedUser := &domain.User{
		ID:        1,
		Username:  username,
		Balance:   1000,
		CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	mock.ExpectQuery(`INSERT INTO users \(username, balance\) VALUES \(\$1, 1000\) RETURNING id, username, balance, created_at`).
		WithArgs(username).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "balance", "created_at"}).
			AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Balance, expectedUser.CreatedAt))

	repo := userRepo.NewUserRepository(db)
	useCase := user.NewUseCase(repo)

	newUser, err := useCase.Register(context.Background(), username)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, newUser)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// Тест на ошибку при регистрации пользователя.
func TestUseCase_Register_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	username := "alice"

	mock.ExpectQuery(`INSERT INTO users \(username, balance\) VALUES \(\$1, 1000\) RETURNING id, username, balance, created_at`).
		WithArgs(username).
		WillReturnError(errors.New("database error"))

	repo := userRepo.NewUserRepository(db)
	useCase := user.NewUseCase(repo)

	newUser, err := useCase.Register(context.Background(), username)
	assert.Error(t, err)
	assert.Nil(t, newUser)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
