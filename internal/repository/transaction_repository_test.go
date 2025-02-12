package repository_test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KazikovAP/merch_store/internal/model"
	"github.com/KazikovAP/merch_store/internal/repository"
)

func TestCreateTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	txn := &model.Transaction{
		UserID:    1,
		Type:      "sent",
		OtherUser: "bob",
		Amount:    200,
	}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO transactions (user_id, type, other_user, amount) VALUES ($1, $2, $3, $4)")).
		WithArgs(txn.UserID, txn.Type, txn.OtherUser, txn.Amount).
		WillReturnResult(sqlmock.NewResult(1, 1))

	txnRepo := repository.NewTransactionRepository(db)

	err = txnRepo.Create(txn)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestGetTransactionsByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "user_id", "type", "other_user", "amount"}).
		AddRow(1, 1, "sent", "bob", 200).
		AddRow(2, 1, "received", "alice", 100)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, type, other_user, amount FROM transactions WHERE user_id=$1")).
		WithArgs(1).
		WillReturnRows(rows)

	txnRepo := repository.NewTransactionRepository(db)

	txns, err := txnRepo.GetByUserID(1)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if len(txns) != 2 {
		t.Errorf("expected 2 transactions, got %d", len(txns))
	}
}
