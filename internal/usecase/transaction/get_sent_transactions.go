package transaction

import (
	"context"
	"fmt"

	"github.com/KazikovAP/merch_store/internal/domain"
)

func (u *useCase) GetSentTransactions(ctx context.Context, userID int) ([]domain.Transaction, error) {
	transactions, err := u.transactionRepo.GetBySenderID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sent transactions for user %d: %w", userID, err)
	}

	return transactions, nil
}
