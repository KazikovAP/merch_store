package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/KazikovAP/merch_store/internal/api/http/auth"
	"github.com/KazikovAP/merch_store/internal/api/http/handlers"
	"github.com/KazikovAP/merch_store/internal/api/http/router"
	"github.com/KazikovAP/merch_store/internal/config"
	"github.com/KazikovAP/merch_store/internal/repository"
	"github.com/KazikovAP/merch_store/internal/usecase"
)

func main() {
	// Загрузка конфигурации
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Инициализация БД
	db, err := initializeDatabase(cfg.DB.DSN())
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	// Инициализация репозиториев
	repo := repository.NewRepositories(db)

	// Инициализация JWT manager
	tokenManager, err := auth.NewJWTManager(cfg.Auth.SigningKey, cfg.Auth.TokenTTL)
	if err != nil {
		if cerr := db.Close(); cerr != nil {
			log.Printf("Ошибка при закрытии БД: %v", cerr)
		}

		log.Fatalf("failed to initialize token manager: %v", err)
	}

	defer db.Close()

	// Инициализация use cases
	useCases := usecase.NewUseCases(repo)

	// Инициализация хендлеров
	handler := handlers.NewHandler(
		useCases.User,
		useCases.Transaction,
		useCases.Purchase,
		useCases.Merch,
		tokenManager,
	)

	// Инициализация роутера
	httpRouter := router.NewRouter(handler, tokenManager)

	// Запуск HTTP сервера
	startServer(httpRouter, cfg.Server.Port, cfg.Server.ReadTimeout, cfg.Server.WriteTimeout)
}

// loadConfig загружает конфигурацию.
func loadConfig() (*config.Config, error) {
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return cfg, nil
}

// initializeDatabase инициализирует соединение с базой данных.
func initializeDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Проверка соединения
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// startServer запускает HTTP сервер с graceful shutdown.
func startServer(r http.Handler, port int, readTimeout, writeTimeout time.Duration) {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	// Graceful shutdown
	go func() {
		log.Println("Server is starting")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := srv.Shutdown(ctx); err != nil {
		cancel()
		log.Fatalf("Server Shutdown: %v", err)
	}

	cancel()

	log.Println("Server exiting")
}
