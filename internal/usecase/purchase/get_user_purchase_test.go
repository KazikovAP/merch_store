package purchase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KazikovAP/merch_store/internal/domain"
	purRep "github.com/KazikovAP/merch_store/internal/repository/purchase"
	"github.com/KazikovAP/merch_store/internal/repository/user"
	"github.com/KazikovAP/merch_store/internal/usecase/purchase"
	"github.com/stretchr/testify/assert"
)

// Тест на успешное получение покупок пользователя.
func TestUseCase_GetUserPurchases_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	userID := 1
	expectedPurchases := []domain.Purchase{
		{
			ID:         1,
			UserID:     userID,
			MerchName:  "t-shirt",
			Quantity:   2,
			TotalPrice: 160,
			CreatedAt:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	mock.ExpectQuery(`SELECT id, username, balance, created_at FROM users WHERE id = \$1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "balance", "created_at"}).
			AddRow(userID, "alice", 1000, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))

	mock.ExpectQuery(`SELECT id, user_id, merch_name, quantity, total_price, created_at FROM purchases WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "merch_name", "quantity", "total_price", "created_at"}).
			AddRow(1, userID, "t-shirt", 2, 160, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))

	repo := purRep.NewPurchaseRepository(db)
	userRepo := user.NewUserRepository(db)
	useCase := purchase.NewUseCase(repo, userRepo, nil)

	purchases, err := useCase.GetUserPurchases(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedPurchases, purchases)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// Тест на ошибку при отсутствии пользователя.
func TestUseCase_GetUserPurchases_UserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	userID := 1

	mock.ExpectQuery(`SELECT id, username, balance, created_at FROM users WHERE id = \$1`).
		WithArgs(userID).
		WillReturnError(errors.New("user not found"))

	repo := purRep.NewPurchaseRepository(db)
	userRepo := user.NewUserRepository(db)
	useCase := purchase.NewUseCase(repo, userRepo, nil)

	purchases, err := useCase.GetUserPurchases(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, purchases)
}

// Тест на ошибку при получении покупок.
func TestUseCase_GetUserPurchases_PurchaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	userID := 1

	mock.ExpectQuery(`SELECT id, username, balance, created_at FROM users WHERE id = \$1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "balance", "created_at"}).
			AddRow(userID, "alice", 1000, "2023-01-01"))

	mock.ExpectQuery(`SELECT id, user_id, merch_name, quantity, total_price, created_at FROM purchases WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnError(errors.New("database error"))

	repo := purRep.NewPurchaseRepository(db)
	userRepo := user.NewUserRepository(db)
	useCase := purchase.NewUseCase(repo, userRepo, nil)

	purchases, err := useCase.GetUserPurchases(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, purchases)
}
