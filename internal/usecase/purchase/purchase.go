package purchase

import (
	"context"
	"fmt"
)

func (u *useCase) Purchase(ctx context.Context, userID, quantity int, merchName string) error {
	// Проверяем существование пользователя и его баланс
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user %d: %w", userID, err)
	}

	// Проверяем существование товара и его цену
	merch, err := u.merchRepo.GetByName(ctx, merchName)
	if err != nil {
		return fmt.Errorf("failed to get merchandise %s: %w", merchName, err)
	}

	// Валидация
	if quantity <= 0 {
		return fmt.Errorf("invalid quantity: %d", quantity)
	}

	totalPrice := merch.Price * quantity
	if user.Balance < totalPrice {
		return fmt.Errorf("insufficient funds: have %d, need %d", user.Balance, totalPrice)
	}

	// Атомарная операция покупки
	if err := u.purchaseRepo.Purchase(ctx, userID, merchName, quantity); err != nil {
		return fmt.Errorf("failed to process purchase: %w", err)
	}

	return nil
}
