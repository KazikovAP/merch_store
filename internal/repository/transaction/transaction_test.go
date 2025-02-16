package transaction_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KazikovAP/merch_store/internal/domain"
	"github.com/KazikovAP/merch_store/internal/repository/transaction"
	"github.com/stretchr/testify/assert"
)

// Тест на успешное создание транзакции.
func TestRepo_Create_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	senderID := 1
	receiverID := 2
	amount := 50

	mock.ExpectBegin()

	mock.ExpectExec(`UPDATE users SET balance = balance - \$1 WHERE id = \$2 AND balance >= \$1`).
		WithArgs(amount, senderID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(`UPDATE users SET balance = balance \+ \$1 WHERE id = \$2`).
		WithArgs(amount, receiverID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(`INSERT INTO transactions \(sender_id, receiver_id, amount\)`).
		WithArgs(senderID, receiverID, amount).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	repo := transaction.NewTransactionRepository(db)

	err = repo.Create(context.Background(), senderID, receiverID, amount)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// Тест на недостаток средств у отправителя.
func TestRepo_Create_InsufficientFunds(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	senderID := 1
	receiverID := 2
	amount := 50

	mock.ExpectBegin()

	mock.ExpectExec(`UPDATE users SET balance = balance - \$1 WHERE id = \$2 AND balance >= \$1`).
		WithArgs(amount, senderID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectRollback()

	repo := transaction.NewTransactionRepository(db)

	err = repo.Create(context.Background(), senderID, receiverID, amount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient funds")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// Тест на успешное получение транзакций по ID пользователя.
func TestRepo_GetByUserID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	userID := 1
	now := time.Now()
	expectedTransactions := []domain.Transaction{
		{
			ID:           1,
			SenderID:     userID,
			ReceiverID:   2,
			Amount:       50,
			CreatedAt:    now,
			SenderName:   "alice",
			ReceiverName: "bob",
		},
		{
			ID:           2,
			SenderID:     3,
			ReceiverID:   userID,
			Amount:       20,
			CreatedAt:    now.Add(-time.Hour),
			SenderName:   "charlie",
			ReceiverName: "alice",
		},
	}

	rows := sqlmock.NewRows([]string{"id", "sender_id", "receiver_id", "amount", "created_at", "sender_name", "receiver_name"})
	for _, t := range expectedTransactions {
		rows.AddRow(t.ID, t.SenderID, t.ReceiverID, t.Amount, t.CreatedAt, t.SenderName, t.ReceiverName)
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
		WillReturnRows(rows)

	repo := transaction.NewTransactionRepository(db)

	transactions, err := repo.GetByUserID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedTransactions, transactions)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// Тест на ошибку при получении транзакций по ID пользователя.
func TestRepo_GetByUserID_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	userID := 1

	mock.ExpectQuery(`SELECT t.id,
	t.sender_id,
	t.receiver_id,
	t.amount,
	t.created_at,
	s.username as sender_name,
	r.username as receiver_name FROM transactions t JOIN users s ON t.sender_id =
	 s.id JOIN users r ON t.receiver_id = r.id WHERE t.sender_id = \$1 OR t.receiver_id = \$1`).
		WithArgs(userID).
		WillReturnError(errors.New("database error"))

	repo := transaction.NewTransactionRepository(db)

	transactions, err := repo.GetByUserID(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, transactions)
	assert.Contains(t, err.Error(), "query transactions")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// Аналогичные тесты для GetBySenderID и GetByReceiverID.
func TestRepo_GetBySenderID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	senderID := 1
	now := time.Now()
	expectedTransactions := []domain.Transaction{
		{
			ID:           1,
			SenderID:     senderID,
			ReceiverID:   2,
			Amount:       50,
			CreatedAt:    now,
			SenderName:   "alice",
			ReceiverName: "bob",
		},
	}

	rows := sqlmock.NewRows([]string{"id", "sender_id", "receiver_id", "amount", "created_at", "sender_name", "receiver_name"})
	for _, t := range expectedTransactions {
		rows.AddRow(t.ID, t.SenderID, t.ReceiverID, t.Amount, t.CreatedAt, t.SenderName, t.ReceiverName)
	}

	mock.ExpectQuery(`SELECT t.id,
	t.sender_id,
	t.receiver_id,
	t.amount,
	t.created_at,
	s.username as sender_name,
	r.username as receiver_name FROM transactions t JOIN users s ON t.sender_id =
	 s.id JOIN users r ON t.receiver_id = r.id WHERE t.sender_id = \$1`).
		WithArgs(senderID).
		WillReturnRows(rows)

	repo := transaction.NewTransactionRepository(db)

	transactions, err := repo.GetBySenderID(context.Background(), senderID)
	assert.NoError(t, err)
	assert.Equal(t, expectedTransactions, transactions)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestRepo_GetByReceiverID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	receiverID := 1
	now := time.Now()
	expectedTransactions := []domain.Transaction{
		{
			ID:           1,
			SenderID:     2,
			ReceiverID:   receiverID,
			Amount:       50,
			CreatedAt:    now,
			SenderName:   "bob",
			ReceiverName: "alice",
		},
	}

	rows := sqlmock.NewRows([]string{"id", "sender_id", "receiver_id", "amount", "created_at", "sender_name", "receiver_name"})
	for _, t := range expectedTransactions {
		rows.AddRow(t.ID, t.SenderID, t.ReceiverID, t.Amount, t.CreatedAt, t.SenderName, t.ReceiverName)
	}

	mock.ExpectQuery(`SELECT t.id,
	t.sender_id,
	t.receiver_id,
	t.amount,
	t.created_at,
	s.username as sender_name,
	r.username as receiver_name FROM transactions t JOIN users s ON t.sender_id =
	 s.id JOIN users r ON t.receiver_id = r.id WHERE t.receiver_id = \$1`).
		WithArgs(receiverID).
		WillReturnRows(rows)

	repo := transaction.NewTransactionRepository(db)

	transactions, err := repo.GetByReceiverID(context.Background(), receiverID)
	assert.NoError(t, err)
	assert.Equal(t, expectedTransactions, transactions)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
