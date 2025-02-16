package merch

import (
	"context"
	"fmt"

	"github.com/KazikovAP/merch_store/internal/domain"
)

func (u *useCase) List(ctx context.Context) ([]domain.Merchandise, error) {
	merch, err := u.merchRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get merchandise list: %w", err)
	}

	return merch, nil
}
