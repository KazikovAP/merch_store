package merch

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/KazikovAP/merch_store/internal/domain"
)

type Repository interface {
	GetByName(ctx context.Context, name string) (*domain.Merchandise, error)
	List(ctx context.Context) ([]domain.Merchandise, error)
}

type Repo struct {
	db *sql.DB
}

func NewMerchRepository(db *sql.DB) Repository {
	return &Repo{db: db}
}

func (r *Repo) GetByName(ctx context.Context, name string) (*domain.Merchandise, error) {
	const query = `
       SELECT name, price
       FROM merchandise
       WHERE name = $1`

	var merch domain.Merchandise
	err := r.db.QueryRowContext(ctx, query, name).
		Scan(&merch.Name, &merch.Price)

	if err != nil {
		return nil, fmt.Errorf("failed to get merchandise by name: %w", err)
	}

	return &merch, nil
}

func (r *Repo) List(ctx context.Context) ([]domain.Merchandise, error) {
	const query = `
       SELECT name, price
       FROM merchandise`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query merchandise: %w", err)
	}
	defer rows.Close()

	var items []domain.Merchandise

	for rows.Next() {
		var item domain.Merchandise

		if err := rows.Scan(&item.Name, &item.Price); err != nil {
			return nil, fmt.Errorf("failed to scan merchandise: %w", err)
		}

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning merchandise: %w", err)
	}

	return items, nil
}
