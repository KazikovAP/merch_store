package service

import (
	"errors"

	"github.com/KazikovAP/merch_store/internal/model"
	"github.com/KazikovAP/merch_store/internal/repository"
)

type CoinService interface {
	TransferCoins(fromUsername, toUsername string, amount int) error
}

type coinService struct {
	userRepo        repository.UserRepository
	transactionRepo repository.TransactionRepository
}

func NewCoinService(userRepo repository.UserRepository, transactionRepo repository.TransactionRepository) CoinService {
	return &coinService{
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *coinService) TransferCoins(fromUsername, toUsername string, amount int) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	sender, err := s.userRepo.GetByUsername(fromUsername)
	if err != nil {
		return err
	}

	if sender.Coins < amount {
		return errors.New("insufficient funds")
	}

	receiver, err := s.userRepo.GetByUsername(toUsername)
	if err != nil {
		return errors.New("receiver not found")
	}

	sender.Coins -= amount
	receiver.Coins += amount

	if err := s.userRepo.Update(sender); err != nil {
		return err
	}

	if err := s.userRepo.Update(receiver); err != nil {
		return err
	}

	if err := s.transactionRepo.Create(&model.Transaction{
		UserID:    sender.ID,
		Type:      "sent",
		OtherUser: receiver.Username,
		Amount:    amount,
	}); err != nil {
		return err
	}

	if err := s.transactionRepo.Create(&model.Transaction{
		UserID:    receiver.ID,
		Type:      "received",
		OtherUser: sender.Username,
		Amount:    amount,
	}); err != nil {
		return err
	}

	return nil
}
