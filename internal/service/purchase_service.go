package service

import (
	"errors"

	"github.com/KazikovAP/merch_store/internal/repository"
)

type PurchaseService interface {
	PurchaseItem(username, item string) error
}

type purchaseService struct {
	userRepo      repository.UserRepository
	inventoryRepo repository.InventoryRepository
	merchRepo     repository.MerchRepository
}

func NewPurchaseService(
	userRepo repository.UserRepository,
	inventoryRepo repository.InventoryRepository,
	merchRepo repository.MerchRepository,
) PurchaseService {
	return &purchaseService{
		userRepo:      userRepo,
		inventoryRepo: inventoryRepo,
		merchRepo:     merchRepo,
	}
}

func (s *purchaseService) PurchaseItem(username, item string) error {
	price, err := s.merchRepo.GetPriceByName(item)
	if err != nil {
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
