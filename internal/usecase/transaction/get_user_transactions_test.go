package transaction_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KazikovAP/merch_store/internal/domain"
	tranRep "github.com/KazikovAP/merch_store/internal/repository/transaction"
	"github.com/KazikovAP/merch_store/internal/usecase/transaction"
	"github.com/stretchr/testify/assert"
)

// Тест на успешное получение всех транзакций пользователя.
func TestUseCase_GetUserTransactions_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	userID := 1
	expectedTransactions := []domain.Transaction{
		{
			ID:           1,
			SenderID:     userID,
			ReceiverID:   2,
			SenderName:   "Alice",
			ReceiverName: "Bob",
			Amount:       100,
			CreatedAt:    time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:           2,
			SenderID:     3,
			ReceiverID:   userID,
			SenderName:   "Charlie",
			ReceiverName: "Alice",
			Amount:       50,
			CreatedAt:    time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
		},
	}

	mock.ExpectQuery(`SELECT t.id,
	t.sender_id,
	t.receiver_id,
	t.amount,
	t.created_at,
	s.username as sender_name,
	r.username as receiver_name FROM transactions t JOIN users s ON t.sender_id =
	 s.id JOIN users r ON t.receiver_id = r.id WHERE t.sender_id = \$1 OR t.receiver_id = \$1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "sender_id", "receiver_id", "amount", "created_at", "sender_name", "receiver_name"}).
			AddRow(1, userID, 2, 100, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), "Alice", "Bob").
			AddRow(2, 3, userID, 50, time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), "Charlie", "Alice"))

	repo := tranRep.NewTransactionRepository(db)
	useCase := transaction.NewUseCase(repo, nil)

	transactions, err := useCase.GetUserTransactions(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedTransactions, transactions)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// Тест на ошибку при получении всех транзакций пользователя.
func TestUseCase_GetUserTransactions_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	userID := 1

	mock.ExpectQuery(`SELECT id, sender_id, receiver_id, amount, created_at FROM transactions WHERE sender_id = \$1 OR receiver_id = \$1`).
		WithArgs(userID, userID).
		WillReturnError(errors.New("database error"))

	repo := tranRep.NewTransactionRepository(db)
	useCase := transaction.NewUseCase(repo, nil)

	transactions, err := useCase.GetUserTransactions(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, transactions)
}
