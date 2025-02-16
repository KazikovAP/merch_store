package purchase

import (
	"context"

	"github.com/KazikovAP/merch_store/internal/domain"
	"github.com/KazikovAP/merch_store/internal/repository/merch"
	"github.com/KazikovAP/merch_store/internal/repository/purchase"
	"github.com/KazikovAP/merch_store/internal/repository/user"
)

type UseCase interface {
	Purchase(ctx context.Context, userID, quantity int, merchName string) error
	GetUserPurchases(ctx context.Context, userID int) ([]domain.Purchase, error)
}

type useCase struct {
	purchaseRepo purchase.Repository
	userRepo     user.Repository
	merchRepo    merch.Repository
}

func NewUseCase(purchaseRepo purchase.Repository, userRepo user.Repository, merchRepo merch.Repository) UseCase {
	return &useCase{
		purchaseRepo: purchaseRepo,
		userRepo:     userRepo,
		merchRepo:    merchRepo,
	}
}
