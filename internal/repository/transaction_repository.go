package repository

import (
	"database/sql"

	"github.com/KazikovAP/merch_store/internal/model"
)

type TransactionRepository interface {
	Create(txn *model.Transaction) error
	GetByUserID(userID int) ([]*model.Transaction, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(txn *model.Transaction) error {
	_, err := r.db.Exec("INSERT INTO transactions (user_id, type, other_user, amount) VALUES ($1, $2, $3, $4)",
		txn.UserID, txn.Type, txn.OtherUser, txn.Amount)
	return err
}

func (r *transactionRepository) GetByUserID(userID int) ([]*model.Transaction, error) {
	rows, err := r.db.Query("SELECT id, user_id, type, other_user, amount FROM transactions WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txns []*model.Transaction

	for rows.Next() {
		txn := &model.Transaction{}
		if err := rows.Scan(&txn.ID, &txn.UserID, &txn.Type, &txn.OtherUser, &txn.Amount); err != nil {
			return nil, err
		}

		txns = append(txns, txn)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return txns, nil
}
