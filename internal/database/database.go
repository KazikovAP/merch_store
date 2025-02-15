package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/KazikovAP/merch_store/internal/config"
)

func SetupDatabase(cfg *config.DatabaseConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	retryCount := cfg.RetryCount
	if retryCount == 0 {
		retryCount = 5
	}

	retryPause := cfg.RetryPause
	if retryPause == 0 {
		retryPause = 2 * time.Second
	}

	var pingErr error
	for i := 0; i < retryCount; i++ {
		pingErr = db.Ping()

		if pingErr == nil {
			break
		}

		log.Printf("Ошибка проверки подключения к БД: %v. Повтор через %v...", pingErr, retryPause)
		time.Sleep(retryPause)
	}

	if pingErr != nil {
		return nil, pingErr
	}

	return db, nil
}
