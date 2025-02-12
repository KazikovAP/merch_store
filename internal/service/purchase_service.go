package service

import (
	"errors"

	"github.com/KazikovAP/merch_store/internal/repository"
)

var MerchItems = map[string]int{
	"t-shirt":    80,
	"cup":        20,
	"book":       50,
	"pen":        10,
	"powerbank":  200,
	"hoody":      300,
	"umbrella":   200,
	"socks":      10,
	"wallet":     50,
	"pink-hoody": 500,
}

type PurchaseService interface {
	PurchaseItem(username, item string) error
}

type purchaseService struct {
	userRepo      repository.UserRepository
	inventoryRepo repository.InventoryRepository
}

func NewPurchaseService(userRepo repository.UserRepository, inventoryRepo repository.InventoryRepository) PurchaseService {
	return &purchaseService{
		userRepo:      userRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (s *purchaseService) PurchaseItem(username, item string) error {
	price, ok := MerchItems[item]
	if !ok {
		return errors.New("item not found")
	}

	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return err
	}

	if user.Coins < price {
		return errors.New("insufficient funds to purchase item")
	}

	user.Coins -= price
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	return s.inventoryRepo.AddItem(user.ID, item, 1)
}
