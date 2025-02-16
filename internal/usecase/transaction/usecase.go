package transaction

import (
	"context"

	"github.com/KazikovAP/merch_store/internal/domain"
	"github.com/KazikovAP/merch_store/internal/repository/transaction"
	"github.com/KazikovAP/merch_store/internal/repository/user"
)

type UseCase interface {
	Transfer(ctx context.Context, senderID, receiverID int, amount int) error
	GetUserTransactions(ctx context.Context, userID int) ([]domain.Transaction, error)
	GetReceivedTransactions(ctx context.Context, userID int) ([]domain.Transaction, error)
	GetSentTransactions(ctx context.Context, userID int) ([]domain.Transaction, error)
}

type useCase struct {
	transactionRepo transaction.Repository
	userRepo        user.Repository
}

func NewUseCase(transactionRepo transaction.Repository, userRepo user.Repository) UseCase {
	return &useCase{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
	}
}
