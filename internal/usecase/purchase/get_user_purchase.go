package purchase

import (
	"context"
	"fmt"

	"github.com/KazikovAP/merch_store/internal/domain"
)

func (u *useCase) GetUserPurchases(ctx context.Context, userID int) ([]domain.Purchase, error) {
	_, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user %d: %w", userID, err)
	}

	purchases, err := u.purchaseRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user purchases: %w", err)
	}

	return purchases, nil
}
