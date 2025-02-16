package transaction

import (
	"context"
	"fmt"
)

func (u *useCase) Transfer(ctx context.Context, senderID, receiverID, amount int) error {
	// Проверяем пользователей
	sender, err := u.userRepo.GetByID(ctx, senderID)
	if err != nil {
		return fmt.Errorf("failed to get sender %d: %w", senderID, err)
	}

	_, err = u.userRepo.GetByID(ctx, receiverID)
	if err != nil {
		return fmt.Errorf("failed to get receiver %d: %w", receiverID, err)
	}

	if senderID == receiverID {
		return fmt.Errorf("sender and receiver are the same user: %d", senderID)
	}

	if amount <= 0 {
		return fmt.Errorf("invalid amount: %d", amount)
	}

	if sender.Balance < amount {
		return fmt.Errorf("insufficient funds: have %d, need %d", sender.Balance, amount)
	}

	// Атомарная операция
	if err := u.transactionRepo.Create(ctx, senderID, receiverID, amount); err != nil {
		return fmt.Errorf("failed to transfer money: %w", err)
	}

	return nil
}
