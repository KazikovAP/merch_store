package repository

import (
	"database/sql"

	"github.com/KazikovAP/merch_store/internal/model"
)

type UserRepository interface {
	GetByUsername(username string) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	user := &model.User{}

	err := r.db.QueryRow("SELECT id, username, password, coins FROM users WHERE username=$1", username).
		Scan(&user.ID, &user.Username, &user.Password, &user.Coins)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.QueryRow(
		"INSERT INTO users (username, password, coins) VALUES ($1, $2, $3) RETURNING id",
		user.Username, user.Password, user.Coins,
	).Scan(&user.ID)
}

func (r *userRepository) Update(user *model.User) error {
	_, err := r.db.Exec("UPDATE users SET password=$1, coins=$2 WHERE id=$3",
		user.Password, user.Coins, user.ID)
	return err
}
