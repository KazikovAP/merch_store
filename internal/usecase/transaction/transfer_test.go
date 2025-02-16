package transaction_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	tranRep "github.com/KazikovAP/merch_store/internal/repository/transaction"
	"github.com/KazikovAP/merch_store/internal/repository/user"
	"github.com/KazikovAP/merch_store/internal/usecase/transaction"
	"github.com/stretchr/testify/assert"
)

// Тест на ошибку при отсутствии отправителя.
func TestUseCase_Transfer_SenderNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	senderID := 1
	receiverID := 2
	amount := 100

	mock.ExpectQuery(`SELECT id, username, balance, created_at FROM users WHERE id = \$1`).
		WithArgs(senderID).
		WillReturnError(errors.New("sender not found"))

	repo := tranRep.NewTransactionRepository(db)
	userRepo := user.NewUserRepository(db)
	useCase := transaction.NewUseCase(repo, userRepo)

	err = useCase.Transfer(context.Background(), senderID, receiverID, amount)
	assert.Error(t, err)
}

// Тест на недостаток средств у отправителя.
func TestUseCase_Transfer_InsufficientFunds(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	senderID := 1
	receiverID := 2
	amount := 100

	mock.ExpectQuery(`SELECT id, username, balance, created_at FROM users WHERE id = \$1`).
		WithArgs(senderID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "balance", "created_at"}).
			AddRow(senderID, "alice", 50, "2023-01-01"))

	mock.ExpectQuery(`SELECT id, username, balance, created_at FROM users WHERE id = \$1`).
		WithArgs(receiverID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "balance", "created_at"}).
			AddRow(receiverID, "bob", 50, "2023-01-01"))

	repo := tranRep.NewTransactionRepository(db)
	userRepo := user.NewUserRepository(db)
	useCase := transaction.NewUseCase(repo, userRepo)

	err = useCase.Transfer(context.Background(), senderID, receiverID, amount)
	assert.Error(t, err)
}
