package user

import (
	"context"

	"github.com/KazikovAP/merch_store/internal/domain"
	"github.com/KazikovAP/merch_store/internal/repository/user"
)

type UseCase interface {
	Register(ctx context.Context, username string) (*domain.User, error)
	GetByID(ctx context.Context, id int) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
}

type useCase struct {
	userRepo user.Repository
}

func NewUseCase(userRepo user.Repository) UseCase {
	return &useCase{
		userRepo: userRepo,
	}
}
