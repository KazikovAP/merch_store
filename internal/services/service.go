package services

import (
	"database/sql"

	"github.com/KazikovAP/merch_store/internal/config"
	"github.com/KazikovAP/merch_store/internal/handlers"
	"github.com/KazikovAP/merch_store/internal/repository"
	"github.com/KazikovAP/merch_store/internal/service"
)

func SetupServicesAndHandlers(db *sql.DB, authCfg config.AuthConfig) *handlers.Handler {
	userRepo := repository.NewUserRepository(db)
	txnRepo := repository.NewTransactionRepository(db)
	invRepo := repository.NewInventoryRepository(db)
	merchRepo := repository.NewMerchRepository(db)

	authService := service.NewAuthService(userRepo, authCfg.JWTSecret)
	userService := service.NewUserService(userRepo, txnRepo, invRepo)
	coinService := service.NewCoinService(userRepo, txnRepo)
	purchaseService := service.NewPurchaseService(userRepo, invRepo, merchRepo)

	return handlers.NewHandler(authService, userService, coinService, purchaseService)
}
