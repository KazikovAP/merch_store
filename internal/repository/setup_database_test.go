package repository_test

import (
	"testing"
	"time"

	"github.com/KazikovAP/merch_store/internal/config"
	"github.com/KazikovAP/merch_store/internal/repository"
	"github.com/stretchr/testify/assert"
)

// Тест на ошибку при открытии соединения.
func TestNewPostgresDB_OpenError(t *testing.T) {
	cfg := &config.DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		Username: "user",
		Password: "pass",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

	pgCfg := repository.PostgresConfig{
		RetryCount: 1,
		RetryPause: time.Millisecond,
	}

	dbConn, err := repository.NewPostgresDB(cfg, pgCfg)
	assert.Error(t, err)
	assert.Nil(t, dbConn)
	assert.Contains(t, err.Error(), "failed to open database")
}
