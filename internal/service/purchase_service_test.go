package service_test

import (
	"testing"

	"github.com/KazikovAP/merch_store/internal/model"
	"github.com/KazikovAP/merch_store/internal/service"
)

func TestPurchaseItem_Success(t *testing.T) {
	user := &model.User{
		ID:       1,
		Username: "testuser",
		Coins:    1000,
	}

	userRepo := &mockUserRepo{user: user}
	invRepo := &mockInventoryRepo{items: []*model.InventoryItem{}}
	purchaseService := service.NewPurchaseService(userRepo, invRepo)

	if err := purchaseService.PurchaseItem("testuser", "t-shirt"); err != nil {
		t.Errorf("expected success, got error: %v", err)
	}

	if user.Coins != 1000-service.MerchItems["t-shirt"] {
		t.Errorf("expected coins to be deducted properly, got %d", user.Coins)
	}

	items, _ := invRepo.GetByUserID(user.ID)
	if len(items) != 1 || items[0].Quantity != 1 {
		t.Errorf("expected inventory to have one item, got %v", items)
	}
}

func TestPurchaseItem_InsufficientFunds(t *testing.T) {
	user := &model.User{
		ID:       1,
		Username: "testuser",
		Coins:    50,
	}

	userRepo := &mockUserRepo{user: user}
	invRepo := &mockInventoryRepo{items: []*model.InventoryItem{}}
	purchaseService := service.NewPurchaseService(userRepo, invRepo)

	if err := purchaseService.PurchaseItem("testuser", "hoody"); err == nil {
		t.Errorf("expected error due to insufficient funds")
	}
}
