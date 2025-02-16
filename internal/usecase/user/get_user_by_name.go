package user

import (
	"context"
	"fmt"

	"github.com/KazikovAP/merch_store/internal/domain"
)

func (u *useCase) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username %s: %w", username, err)
	}

	return user, nil
}
