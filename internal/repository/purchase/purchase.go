package purchase

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/KazikovAP/merch_store/internal/domain"
)

type Repository interface {
	Purchase(ctx context.Context, userID int, merchName string, quantity int) error
	GetByUserID(ctx context.Context, userID int) ([]domain.Purchase, error)
}

type Repo struct {
	db *sql.DB
}

func NewPurchaseRepository(db *sql.DB) Repository {
	return &Repo{db: db}
}

func (r *Repo) Purchase(ctx context.Context, userID int, merchName string, quantity int) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			fmt.Printf("rollback failed: %v\n", err)
		}
	}()

	// Получаем цену товара
	var price int

	err = tx.QueryRowContext(ctx, `
       SELECT price 
       FROM merchandise 
       WHERE name = $1`, merchName).Scan(&price)

	if err != nil {
		return fmt.Errorf("get merchandise price: %w", err)
	}

	totalPrice := price * quantity

	// Списываем деньги с баланса пользователя
	result, err := tx.ExecContext(ctx, `
       UPDATE users 
       SET balance = balance - $1 
       WHERE id = $2 AND balance >= $1`, totalPrice, userID)
	if err != nil {
		return fmt.Errorf("update user balance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("insufficient funds")
	}

	// Создаем запись о покупке
	_, err = tx.ExecContext(ctx, `
       INSERT INTO purchases (user_id, merch_name, quantity, total_price)
       VALUES ($1, $2, $3, $4)`, userID, merchName, quantity, totalPrice)

	if err != nil {
		return fmt.Errorf("create purchase record: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (r *Repo) GetByUserID(ctx context.Context, userID int) ([]domain.Purchase, error) {
	rows, err := r.db.QueryContext(ctx, `
       SELECT id, user_id, merch_name, quantity, total_price, created_at
       FROM purchases
       WHERE user_id = $1
       ORDER BY created_at DESC`, userID)

	if err != nil {
		return nil, fmt.Errorf("query purchases: %w", err)
	}

	defer rows.Close()

	var purchases []domain.Purchase

	for rows.Next() {
		var p domain.Purchase

		if err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.MerchName,
			&p.Quantity,
			&p.TotalPrice,
			&p.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan purchase: %w", err)
		}

		purchases = append(purchases, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning purchases: %w", err)
	}

	return purchases, nil
}
