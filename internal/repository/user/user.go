package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/KazikovAP/merch_store/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, username string) (*domain.User, error)
	GetByID(ctx context.Context, id int) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
}

type Repo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) Repository {
	return &Repo{db: db}
}

func (r *Repo) Create(ctx context.Context, username string) (*domain.User, error) {
	const query = `
        INSERT INTO users (username, balance)
        VALUES ($1, 1000)
        RETURNING id, username, balance, created_at`

	var user domain.User

	err := r.db.QueryRowContext(ctx, query, username).
		Scan(&user.ID, &user.Username, &user.Balance, &user.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (r *Repo) GetByID(ctx context.Context, id int) (*domain.User, error) {
	const query = `
        SELECT id, username, balance, created_at
        FROM users
        WHERE id = $1`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&user.ID, &user.Username, &user.Balance, &user.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

func (r *Repo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	const query = `
        SELECT id, username, balance, created_at
        FROM users
        WHERE username = $1`

	var user domain.User

	err := r.db.QueryRowContext(ctx, query, username).
		Scan(&user.ID, &user.Username, &user.Balance, &user.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &user, nil
}
