package merch

import (
	"context"
	"fmt"

	"github.com/KazikovAP/merch_store/internal/domain"
)

func (u *useCase) GetByName(ctx context.Context, name string) (*domain.Merchandise, error) {
	if name == "" {
		return nil, fmt.Errorf("empty merchandise name")
	}

	merch, err := u.merchRepo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get merchandise by name %s: %w", name, err)
	}

	return merch, nil
}
