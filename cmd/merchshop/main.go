package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/lib/pq"

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

	db, err := setupDatabase(cfg)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка при закрытии базы данных: %v", err)
		}
	}()

	handler := setupServicesAndHandlers(db, cfg)

	router := setupRouter(handler, cfg)

	startServer(cfg.ServerPort, router)
}

func setupDatabase(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if pingErr := db.Ping(); pingErr != nil {
		return nil, pingErr
	}

	return db, nil
}

func setupServicesAndHandlers(db *sql.DB, cfg *config.Config) *handlers.Handler {
	userRepo := repository.NewUserRepository(db)
	txnRepo := repository.NewTransactionRepository(db)
	invRepo := repository.NewInventoryRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	userService := service.NewUserService(userRepo, txnRepo, invRepo)
	coinService := service.NewCoinService(userRepo, txnRepo)
	purchaseService := service.NewPurchaseService(userRepo, invRepo)

	return handlers.NewHandler(authService, userService, coinService, purchaseService)
}

func setupRouter(h *handlers.Handler, cfg *config.Config) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/auth", h.Auth).Methods("POST")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.JwtMiddleware(cfg.JWTSecret))
	api.HandleFunc("/info", h.Info).Methods("GET")
	api.HandleFunc("/sendCoin", h.SendCoin).Methods("POST")
	api.HandleFunc("/buy/{item}", h.Buy).Methods("GET")

	return r
}

func startServer(port string, handler http.Handler) {
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Println("Сервер запущен на порту :" + port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("Останавливаем сервер...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Ошибка при остановке сервера: %v", err)
		os.Exit(1)
	}

	log.Println("Сервер успешно остановлен.")
}
