package purchase_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KazikovAP/merch_store/internal/domain"
	"github.com/KazikovAP/merch_store/internal/repository/purchase"
	"github.com/stretchr/testify/assert"
)

// Тест на успешную покупку товара.
func TestRepo_Purchase_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	userID := 1
	merchName := "t-shirt"
	quantity := 2
	price := 80
	totalPrice := price * quantity

	mock.ExpectBegin()

	mock.ExpectQuery(`SELECT price FROM merchandise WHERE name = \$1`).
		WithArgs(merchName).
		WillReturnRows(sqlmock.NewRows([]string{"price"}).AddRow(price))

	mock.ExpectExec(`UPDATE users SET balance = balance - \$1 WHERE id = \$2 AND balance >= \$1`).
		WithArgs(totalPrice, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(`INSERT INTO purchases \(user_id, merch_name, quantity, total_price\)`).
		WithArgs(userID, merchName, quantity, totalPrice).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	repo := purchase.NewPurchaseRepository(db)

	err = repo.Purchase(context.Background(), userID, merchName, quantity)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// Тест на ошибку при получении цены товара.
func TestRepo_Purchase_MerchNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	userID := 1
	merchName := "unknown-item"
	quantity := 2

	mock.ExpectBegin()

	mock.ExpectQuery(`SELECT price FROM merchandise WHERE name = \$1`).
		WithArgs(merchName).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectRollback()

	repo := purchase.NewPurchaseRepository(db)

	err = repo.Purchase(context.Background(), userID, merchName, quantity)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "get merchandise price")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// Тест на недостаток средств у пользователя.
func TestRepo_Purchase_InsufficientFunds(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	userID := 1
	merchName := "t-shirt"
	quantity := 2
	price := 80
	totalPrice := price * quantity

	mock.ExpectBegin()

	mock.ExpectQuery(`SELECT price FROM merchandise WHERE name = \$1`).
		WithArgs(merchName).
		WillReturnRows(sqlmock.NewRows([]string{"price"}).AddRow(price))

	mock.ExpectExec(`UPDATE users SET balance = balance - \$1 WHERE id = \$2 AND balance >= \$1`).
		WithArgs(totalPrice, userID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectRollback()

	repo := purchase.NewPurchaseRepository(db)

	err = repo.Purchase(context.Background(), userID, merchName, quantity)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient funds")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// Тест на успешное получение списка покупок пользователя.
func TestRepo_GetByUserID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	userID := 1
	now := time.Now()
	expectedPurchases := []domain.Purchase{
		{
			ID:         1,
			UserID:     userID,
			MerchName:  "t-shirt",
			Quantity:   2,
			TotalPrice: 160,
			CreatedAt:  now,
		},
		{
			ID:         2,
			UserID:     userID,
			MerchName:  "cup",
			Quantity:   1,
			TotalPrice: 20,
			CreatedAt:  now.Add(-time.Hour),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "merch_name", "quantity", "total_price", "created_at"})
	for _, p := range expectedPurchases {
		rows.AddRow(p.ID, p.UserID, p.MerchName, p.Quantity, p.TotalPrice, p.CreatedAt)
	}

	mock.ExpectQuery(`SELECT id,
	user_id, merch_name,
	quantity,
	total_price,
	created_at FROM purchases WHERE user_id = \$1 ORDER BY created_at DESC`).
		WithArgs(userID).
		WillReturnRows(rows)

	repo := purchase.NewPurchaseRepository(db)

	purchases, err := repo.GetByUserID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedPurchases, purchases)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// Тест на ошибку при получении списка покупок.
func TestRepo_GetByUserID_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	userID := 1

	mock.ExpectQuery(`SELECT id,
	user_id,
	merch_name,
	quantity,
	total_price,
	created_at FROM purchases WHERE user_id = \$1 ORDER BY created_at DESC`).
		WithArgs(userID).
		WillReturnError(errors.New("database error"))

	repo := purchase.NewPurchaseRepository(db)

	purchases, err := repo.GetByUserID(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, purchases)
	assert.Contains(t, err.Error(), "query purchases")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
