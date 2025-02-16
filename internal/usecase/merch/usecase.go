package merch

import (
	"context"

	"github.com/KazikovAP/merch_store/internal/domain"
	"github.com/KazikovAP/merch_store/internal/repository/merch"
)

type UseCase interface {
	List(ctx context.Context) ([]domain.Merchandise, error)
	GetByName(ctx context.Context, name string) (*domain.Merchandise, error)
}

type useCase struct {
	merchRepo merch.Repository
}

func NewUseCase(mr merch.Repository) UseCase {
	return &useCase{
		merchRepo: mr,
	}
}
