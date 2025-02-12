package repository

import (
	"database/sql"

	"github.com/KazikovAP/merch_store/internal/model"
)

type InventoryRepository interface {
	GetByUserID(userID int) ([]*model.InventoryItem, error)
	AddItem(userID int, itemType string, quantity int) error
}

type inventoryRepository struct {
	db *sql.DB
}

func NewInventoryRepository(db *sql.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) GetByUserID(userID int) ([]*model.InventoryItem, error) {
	rows, err := r.db.Query("SELECT id, user_id, item_type, quantity FROM inventory WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*model.InventoryItem

	for rows.Next() {
		item := &model.InventoryItem{}
		if err := rows.Scan(&item.ID, &item.UserID, &item.ItemType, &item.Quantity); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *inventoryRepository) AddItem(userID int, itemType string, quantity int) error {
	var id int

	err := r.db.QueryRow("SELECT id FROM inventory WHERE user_id=$1 AND item_type=$2", userID, itemType).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = r.db.Exec("INSERT INTO inventory (user_id, item_type, quantity) VALUES ($1, $2, $3)", userID, itemType, quantity)
			return err
		}

		return err
	}

	_, err = r.db.Exec("UPDATE inventory SET quantity = quantity + $1 WHERE id=$2", quantity, id)

	return err
}
