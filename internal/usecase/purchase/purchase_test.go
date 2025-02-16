package purchase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KazikovAP/merch_store/internal/repository/merch"
	purRep "github.com/KazikovAP/merch_store/internal/repository/purchase"
	"github.com/KazikovAP/merch_store/internal/repository/user"
	"github.com/KazikovAP/merch_store/internal/usecase/purchase"
	"github.com/stretchr/testify/assert"
)

// Тест на ошибку при отсутствии пользователя.
func TestUseCase_Purchase_UserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	userID := 1
	merchName := "t-shirt"
	quantity := 2

	mock.ExpectQuery(`SELECT id, username, balance, created_at FROM users WHERE id = \$1`).
		WithArgs(userID).
		WillReturnError(errors.New("user not found"))

	repo := purRep.NewPurchaseRepository(db)
	userRepo := user.NewUserRepository(db)
	merchRepo := merch.NewMerchRepository(db)
	useCase := purchase.NewUseCase(repo, userRepo, merchRepo)

	err = useCase.Purchase(context.Background(), userID, quantity, merchName)
	assert.Error(t, err)
}

// Тест на ошибку при отсутствии товара.
func TestUseCase_Purchase_MerchNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	userID := 1
	merchName := "unknown-item"
	quantity := 2

	mock.ExpectQuery(`SELECT id, username, balance, created_at FROM users WHERE id = \$1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "balance", "created_at"}).
			AddRow(userID, "alice", 1000, "2023-01-01"))

	mock.ExpectQuery(`SELECT name, price FROM merchandise WHERE name = \$1`).
		WithArgs(merchName).
		WillReturnError(errors.New("merchandise not found"))

	repo := purRep.NewPurchaseRepository(db)
	userRepo := user.NewUserRepository(db)
	merchRepo := merch.NewMerchRepository(db)
	useCase := purchase.NewUseCase(repo, userRepo, merchRepo)

	err = useCase.Purchase(context.Background(), userID, quantity, merchName)
	assert.Error(t, err)
}

// Тест на недостаток средств у пользователя.
func TestUseCase_Purchase_InsufficientFunds(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	userID := 1
	merchName := "wallet"
	quantity := 2
	price := 80

	mock.ExpectQuery(`SELECT id, username, balance, created_at FROM users WHERE id = \$1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "balance", "created_at"}).
			AddRow(userID, "alice", 50, "2023-01-01"))

	mock.ExpectQuery(`SELECT name, price FROM merchandise WHERE name = \$1`).
		WithArgs(merchName).
		WillReturnRows(sqlmock.NewRows([]string{"name", "price"}).
			AddRow(merchName, price))

	repo := purRep.NewPurchaseRepository(db)
	userRepo := user.NewUserRepository(db)
	merchRepo := merch.NewMerchRepository(db)
	useCase := purchase.NewUseCase(repo, userRepo, merchRepo)

	err = useCase.Purchase(context.Background(), userID, quantity, merchName)
	assert.Error(t, err)
}

// Тест на ошибку при выполнении транзакции.
func TestUseCase_Purchase_TransactionError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	userID := 1
	merchName := "hoody"
	quantity := 2
	price := 80

	mock.ExpectQuery(`SELECT id, username, balance, created_at FROM users WHERE id = \$1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "balance", "created_at"}).
			AddRow(userID, "alice", 1000, "2023-01-01"))

	mock.ExpectQuery(`SELECT name, price FROM merchandise WHERE name = \$1`).
		WithArgs(merchName).
		WillReturnRows(sqlmock.NewRows([]string{"name", "price"}).
			AddRow(merchName, price))

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE users SET balance = balance - \$1 WHERE id = \$2 AND balance >= \$1`).
		WithArgs(price*quantity, userID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	repo := purRep.NewPurchaseRepository(db)
	userRepo := user.NewUserRepository(db)
	merchRepo := merch.NewMerchRepository(db)
	useCase := purchase.NewUseCase(repo, userRepo, merchRepo)

	err = useCase.Purchase(context.Background(), userID, quantity, merchName)
	assert.Error(t, err)
}
