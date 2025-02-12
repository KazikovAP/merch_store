package service_test

import (
	"testing"

	"github.com/KazikovAP/merch_store/internal/model"
	"github.com/KazikovAP/merch_store/internal/service"
)

func TestGetUserInfo(t *testing.T) {
	user := &model.User{
		ID:       1,
		Username: "testuser",
		Coins:    500,
	}

	userRepo := &mockUserRepo{user: user}

	transactions := []*model.Transaction{
		{
			ID:        1,
			UserID:    1,
			Type:      "received",
			OtherUser: "alice",
			Amount:    100,
		},
		{
			ID:        2,
			UserID:    1,
			Type:      "sent",
			OtherUser: "bob",
			Amount:    50,
		},
	}

	txnRepo := &mockTransactionRepo{transactions: transactions}

	inventoryItems := []*model.InventoryItem{
		{
			ID:       1,
			UserID:   1,
			ItemType: "t-shirt",
			Quantity: 1,
		},
		{
			ID:       2,
			UserID:   1,
			ItemType: "pen",
			Quantity: 2,
		},
	}

	invRepo := &mockInventoryRepo{items: inventoryItems}

	userService := service.NewUserService(userRepo, txnRepo, invRepo)

	info, err := userService.GetUserInfo("testuser")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if info.Coins != user.Coins {
		t.Errorf("expected coins %d, got %d", user.Coins, info.Coins)
	}

	if len(info.Inventory) != len(inventoryItems) {
		t.Errorf("expected %d inventory items, got %d", len(inventoryItems), len(info.Inventory))
	}

	if len(info.CoinHistory.Received) != 1 {
		t.Errorf("expected 1 received transaction, got %d", len(info.CoinHistory.Received))
	}

	if len(info.CoinHistory.Sent) != 1 {
		t.Errorf("expected 1 sent transaction, got %d", len(info.CoinHistory.Sent))
	}

	received := info.CoinHistory.Received[0]
	if received.FromUser != "alice" || received.Amount != 100 {
		t.Errorf("unexpected received transaction: %+v", received)
	}

	sent := info.CoinHistory.Sent[0]
	if sent.ToUser != "bob" || sent.Amount != 50 {
		t.Errorf("unexpected sent transaction: %+v", sent)
	}
}
