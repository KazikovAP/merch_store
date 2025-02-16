package user

import (
	"context"
	"fmt"

	"github.com/KazikovAP/merch_store/internal/domain"
)

func (u *useCase) Register(ctx context.Context, username string) (*domain.User, error) {
	user, err := u.userRepo.Create(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}
