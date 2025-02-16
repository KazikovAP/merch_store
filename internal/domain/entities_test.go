package domain_test

import (
	"testing"
	"time"

	"github.com/KazikovAP/merch_store/internal/domain"
	"github.com/stretchr/testify/assert"
)

// Тест на инициализацию структуры User.
func TestUser_Init(t *testing.T) {
	now := time.Now()
	user := domain.User{
		ID:        1,
		Username:  "alice",
		Balance:   1000,
		CreatedAt: now,
	}

	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "alice", user.Username)
	assert.Equal(t, 1000, user.Balance)
	assert.Equal(t, now, user.CreatedAt)
}

// Тест на сравнение двух структур User.
func TestUser_Equal(t *testing.T) {
	now := time.Now()
	user1 := domain.User{
		ID:        1,
		Username:  "alice",
		Balance:   1000,
		CreatedAt: now,
	}
	user2 := domain.User{
		ID:        1,
		Username:  "alice",
		Balance:   1000,
		CreatedAt: now,
	}

	assert.Equal(t, user1, user2)
}

// Тест на инициализацию структуры Merchandise.
func TestMerchandise_Init(t *testing.T) {
	merch := domain.Merchandise{
		Name:  "t-shirt",
		Price: 80,
	}

	assert.Equal(t, "t-shirt", merch.Name)
	assert.Equal(t, 80, merch.Price)
}

// Тест на инициализацию структуры Transaction.
func TestTransaction_Init(t *testing.T) {
	now := time.Now()
	tx := domain.Transaction{
		ID:           1,
		SenderID:     101,
		ReceiverID:   102,
		SenderName:   "alice",
		ReceiverName: "bob",
		Amount:       50,
		CreatedAt:    now,
	}

	assert.Equal(t, 1, tx.ID)
	assert.Equal(t, 101, tx.SenderID)
	assert.Equal(t, 102, tx.ReceiverID)
	assert.Equal(t, "alice", tx.SenderName)
	assert.Equal(t, "bob", tx.ReceiverName)
	assert.Equal(t, 50, tx.Amount)
	assert.Equal(t, now, tx.CreatedAt)
}

// Тест на инициализацию структуры Purchase.
func TestPurchase_Init(t *testing.T) {
	now := time.Now()
	purchase := domain.Purchase{
		ID:         1,
		UserID:     101,
		MerchName:  "t-shirt",
		Quantity:   2,
		TotalPrice: 160,
		CreatedAt:  now,
	}

	assert.Equal(t, 1, purchase.ID)
	assert.Equal(t, 101, purchase.UserID)
	assert.Equal(t, "t-shirt", purchase.MerchName)
	assert.Equal(t, 2, purchase.Quantity)
	assert.Equal(t, 160, purchase.TotalPrice)
	assert.Equal(t, now, purchase.CreatedAt)
}
