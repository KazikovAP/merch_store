package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/KazikovAP/merch_store/internal/config"
)

type PostgresConfig struct {
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	RetryCount      int           `mapstructure:"retry_count"`
	RetryPause      time.Duration `mapstructure:"retry_pause"`
}

func NewPostgresDB(cfg *config.DatabaseConfig, pgCfg PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(pgCfg.MaxOpenConns)
	db.SetMaxIdleConns(pgCfg.MaxIdleConns)
	db.SetConnMaxLifetime(pgCfg.ConnMaxLifetime)

	// Значения по умолчанию для retry
	retryCount := pgCfg.RetryCount
	if retryCount == 0 {
		retryCount = 5
	}

	retryPause := pgCfg.RetryPause
	if retryPause == 0 {
		retryPause = 2 * time.Second
	}

	// Проверка соединения с повторами
	var pingErr error
	for i := 0; i < retryCount; i++ {
		pingErr = db.Ping()
		if pingErr == nil {
			break
		}

		log.Printf("Failed to ping database: %v. Retrying in %v...", pingErr, retryPause)
		time.Sleep(retryPause)
	}

	if pingErr != nil {
		return nil, fmt.Errorf("failed to connect to database after %d retries: %w", retryCount, pingErr)
	}

	return db, nil
}
