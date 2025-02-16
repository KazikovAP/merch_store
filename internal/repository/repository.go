package repository

import (
	"database/sql"

	"github.com/KazikovAP/merch_store/internal/repository/merch"
	"github.com/KazikovAP/merch_store/internal/repository/purchase"
	"github.com/KazikovAP/merch_store/internal/repository/transaction"
	"github.com/KazikovAP/merch_store/internal/repository/user"
)

type Repositories struct {
	User        user.Repository
	Transaction transaction.Repository
	Purchase    purchase.Repository
	Merch       merch.Repository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		User:        user.NewUserRepository(db),
		Transaction: transaction.NewTransactionRepository(db),
		Purchase:    purchase.NewPurchaseRepository(db),
		Merch:       merch.NewMerchRepository(db),
	}
}
