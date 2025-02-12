package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/KazikovAP/merch_store/internal/config"
	"github.com/KazikovAP/merch_store/internal/handlers"
	"github.com/KazikovAP/merch_store/internal/middleware"
	"github.com/KazikovAP/merch_store/internal/repository"
	"github.com/KazikovAP/merch_store/internal/service"
	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Не удалось загрузить конфигурацию:", err)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Ошибка при проверке подключения к БД:", err)
	}

	userRepo := repository.NewUserRepository(db)
	txnRepo := repository.NewTransactionRepository(db)
	invRepo := repository.NewInventoryRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	userService := service.NewUserService(userRepo, txnRepo, invRepo)
	coinService := service.NewCoinService(userRepo, txnRepo)
	purchaseService := service.NewPurchaseService(userRepo, invRepo)

	h := handlers.NewHandler(authService, userService, coinService, purchaseService)

	r := mux.NewRouter()
	r.HandleFunc("/api/auth", h.Auth).Methods("POST")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.JwtMiddleware(cfg.JWTSecret))
	api.HandleFunc("/info", h.Info).Methods("GET")
	api.HandleFunc("/sendCoin", h.SendCoin).Methods("POST")
	api.HandleFunc("/buy/{item}", h.Buy).Methods("GET")

	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("Сервер запущен на :" + cfg.ServerPort)

	if err := srv.ListenAndServe(); err != nil {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Ошибка при закрытии базы данных: %v", closeErr)
		}

		log.Fatal(err)
	}
}
