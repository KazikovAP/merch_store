package transaction

import (
	"context"
	"fmt"

	"github.com/KazikovAP/merch_store/internal/domain"
)

func (u *useCase) GetReceivedTransactions(ctx context.Context, userID int) ([]domain.Transaction, error) {
	transactions, err := u.transactionRepo.GetByReceiverID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get received transactions for user %d: %w", userID, err)
	}

	return transactions, nil
}
