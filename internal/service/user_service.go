package service

import (
	"github.com/KazikovAP/merch_store/internal/model/dto"
	"github.com/KazikovAP/merch_store/internal/repository"
)

type UserService interface {
	GetUserInfo(username string) (*dto.InfoResponse, error)
}

type userService struct {
	userRepo        repository.UserRepository
	transactionRepo repository.TransactionRepository
	inventoryRepo   repository.InventoryRepository
}

func NewUserService(
	userRepo repository.UserRepository,
	transactionRepo repository.TransactionRepository,
	inventoryRepo repository.InventoryRepository,
) UserService {
	return &userService{
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
		inventoryRepo:   inventoryRepo,
	}
}

func (s *userService) GetUserInfo(username string) (*dto.InfoResponse, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	txns, err := s.transactionRepo.GetByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	invItems, err := s.inventoryRepo.GetByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	invResp := make([]dto.InventoryResponse, 0, len(invItems))

	for _, item := range invItems {
		invResp = append(invResp, dto.InventoryResponse{
			Type:     item.ItemType,
			Quantity: item.Quantity,
		})
	}

	coinHistory := dto.CoinHistoryResponse{
		Received: []dto.TransactionDetail{},
		Sent:     []dto.TransactionDetail{},
	}

	for _, txn := range txns {
		if txn.Type == "received" {
			coinHistory.Received = append(coinHistory.Received, dto.TransactionDetail{
				FromUser: txn.OtherUser,
				Amount:   txn.Amount,
			})
		} else if txn.Type == "sent" {
			coinHistory.Sent = append(coinHistory.Sent, dto.TransactionDetail{
				ToUser: txn.OtherUser,
				Amount: txn.Amount,
			})
		}
	}

	info := &dto.InfoResponse{
		Coins:       user.Coins,
		Inventory:   invResp,
		CoinHistory: coinHistory,
	}

	return info, nil
}
