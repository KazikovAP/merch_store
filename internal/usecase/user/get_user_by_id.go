package user

import (
	"context"
	"fmt"

	"github.com/KazikovAP/merch_store/internal/domain"
)

func (u *useCase) GetByID(ctx context.Context, id int) (*domain.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id %d: %w", id, err)
	}

	return user, nil
}
