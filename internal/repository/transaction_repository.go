package repository

import (
	"database/sql"
	"fmt"

	"github.com/KazikovAP/merch_store/internal/model/domain"
)

type TransactionRepository interface {
	Create(txn *domain.Transaction) error
	GetByUserID(userID int) ([]*domain.Transaction, error)
	CreateBatch(transactions []*domain.Transaction) error
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(txn *domain.Transaction) error {
	_, err := r.db.Exec(
		"INSERT INTO transactions (user_id, type, other_user, amount) VALUES ($1, $2, $3, $4)",
		txn.UserID, txn.Type, txn.OtherUser, txn.Amount,
	)

	return err
}

func (r *transactionRepository) GetByUserID(userID int) ([]*domain.Transaction, error) {
	rows, err := r.db.Query(
		"SELECT id, user_id, type, other_user, amount FROM transactions WHERE user_id=$1",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txns []*domain.Transaction

	for rows.Next() {
		txn := &domain.Transaction{}
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

func (r *transactionRepository) CreateBatch(transactions []*domain.Transaction) error {
	if len(transactions) == 0 {
		return nil
	}

	query := "INSERT INTO transactions (user_id, type, other_user, amount) VALUES "

	args := []interface{}{}

	for i, txn := range transactions {
		query += fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4)
		if i < len(transactions)-1 {
			query += ", "
		}

		args = append(args, txn.UserID, txn.Type, txn.OtherUser, txn.Amount)
	}

	_, err := r.db.Exec(query, args...)

	return err
}
