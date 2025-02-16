package usecase

import (
	"github.com/KazikovAP/merch_store/internal/repository"
	"github.com/KazikovAP/merch_store/internal/usecase/merch"
	"github.com/KazikovAP/merch_store/internal/usecase/purchase"
	"github.com/KazikovAP/merch_store/internal/usecase/transaction"
	"github.com/KazikovAP/merch_store/internal/usecase/user"
)

type UseCases struct {
	User        user.UseCase
	Transaction transaction.UseCase
	Purchase    purchase.UseCase
	Merch       merch.UseCase
}

func NewUseCases(repos *repository.Repositories) *UseCases {
	return &UseCases{
		User:        user.NewUseCase(repos.User),
		Transaction: transaction.NewUseCase(repos.Transaction, repos.User),
		Purchase:    purchase.NewUseCase(repos.Purchase, repos.User, repos.Merch),
		Merch:       merch.NewUseCase(repos.Merch),
	}
}
