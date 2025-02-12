package service_test

import (
	"testing"

	"github.com/KazikovAP/merch_store/internal/model"
	"github.com/KazikovAP/merch_store/internal/service"
)

func TestTransferCoins_Success(t *testing.T) {
	userRepo := &mockUserRepoMap{users: map[string]*model.User{
		"alice": {ID: 1, Username: "alice", Coins: 1000},
		"bob":   {ID: 2, Username: "bob", Coins: 1000},
	}}

	txnRepo := &mockTransactionRepo{}
	coinService := service.NewCoinService(userRepo, txnRepo)

	if err := coinService.TransferCoins("alice", "bob", 200); err != nil {
		t.Errorf("expected transfer to succeed, got error: %v", err)
	}

	alice, _ := userRepo.GetByUsername("alice")
	bob, _ := userRepo.GetByUsername("bob")

	if alice.Coins != 800 {
		t.Errorf("expected alice coins to be 800, got %d", alice.Coins)
	}

	if bob.Coins != 1200 {
		t.Errorf("expected bob coins to be 1200, got %d", bob.Coins)
	}
}

func TestTransferCoins_InsufficientFunds(t *testing.T) {
	userRepo := &mockUserRepoMap{users: map[string]*model.User{
		"alice": {ID: 1, Username: "alice", Coins: 100},
		"bob":   {ID: 2, Username: "bob", Coins: 1000},
	}}

	txnRepo := &mockTransactionRepo{}
	coinService := service.NewCoinService(userRepo, txnRepo)

	if err := coinService.TransferCoins("alice", "bob", 200); err == nil {
		t.Errorf("expected error due to insufficient funds")
	}
}
