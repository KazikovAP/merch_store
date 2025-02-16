package transaction

import (
	"context"
	"fmt"

	"github.com/KazikovAP/merch_store/internal/domain"
)

func (u *useCase) GetUserTransactions(ctx context.Context, userID int) ([]domain.Transaction, error) {
	transactions, err := u.transactionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions for user %d: %w", userID, err)
	}

	return transactions, nil
}
