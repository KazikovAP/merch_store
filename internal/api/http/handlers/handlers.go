package handlers

import (
	"github.com/KazikovAP/merch_store/internal/api/http/auth"
	"github.com/KazikovAP/merch_store/internal/usecase/merch"
	"github.com/KazikovAP/merch_store/internal/usecase/purchase"
	"github.com/KazikovAP/merch_store/internal/usecase/transaction"
	"github.com/KazikovAP/merch_store/internal/usecase/user"
)

type Handler struct {
	userUseCase        user.UseCase
	transactionUseCase transaction.UseCase
	purchaseUseCase    purchase.UseCase
	merchUseCase       merch.UseCase
	tokenManager       auth.TokenManager
}

func NewHandler(
	userUseCase user.UseCase,
	transactionUseCase transaction.UseCase,
	purchaseUseCase purchase.UseCase,
	merchUseCase merch.UseCase,
	tm auth.TokenManager,
) *Handler {
	return &Handler{
		userUseCase:        userUseCase,
		transactionUseCase: transactionUseCase,
		purchaseUseCase:    purchaseUseCase,
		merchUseCase:       merchUseCase,
		tokenManager:       tm,
	}
}
