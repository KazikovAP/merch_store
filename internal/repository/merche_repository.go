package repository

import "database/sql"

type MerchRepository interface {
	GetPriceByName(itemName string) (int, error)
}

type merchRepository struct {
	db *sql.DB
}

func NewMerchRepository(db *sql.DB) MerchRepository {
	return &merchRepository{db: db}
}

func (r *merchRepository) GetPriceByName(itemName string) (int, error) {
	var price int

	err := r.db.QueryRow("SELECT price FROM merchandise WHERE name=$1", itemName).Scan(&price)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		return 0, err
	}

	return price, nil
}
