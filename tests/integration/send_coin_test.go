package integration_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KazikovAP/merch_store/internal/api/http/auth"
	"github.com/KazikovAP/merch_store/internal/api/http/handlers"
	"github.com/KazikovAP/merch_store/internal/api/http/router"
	"github.com/KazikovAP/merch_store/internal/repository"
	"github.com/KazikovAP/merch_store/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestSendCoins(t *testing.T) {
	dsn := "host=localhost user=postgres password=postgres sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("DROP DATABASE IF EXISTS testdb")
	if err != nil {
		t.Fatalf("failed to drop test database: %v", err)
	}

	_, err = db.Exec("CREATE DATABASE testdb")
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	testDB, err := sql.Open("postgres", "host=localhost user=postgres password=postgres dbname=testdb sslmode=disable")
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	defer testDB.Close()

	_, err = testDB.Exec(`
        CREATE TABLE users (
            id BIGSERIAL PRIMARY KEY,
            username VARCHAR(50) UNIQUE NOT NULL,
            balance BIGINT NOT NULL CHECK (balance >= 0),
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );
        CREATE TABLE transactions (
            id BIGSERIAL PRIMARY KEY,
            sender_id BIGINT NOT NULL REFERENCES users(id),
            receiver_id BIGINT NOT NULL REFERENCES users(id),
            amount BIGINT NOT NULL CHECK (amount > 0),
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            CONSTRAINT different_users CHECK (sender_id != receiver_id)
        );
    `)

	if err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}

	_, err = testDB.Exec(`
        INSERT INTO users (username, balance) VALUES ('alice', 1000), ('bob', 500);
    `)

	if err != nil {
		t.Fatalf("failed to insert initial data: %v", err)
	}

	repo := repository.NewRepositories(testDB)
	useCases := usecase.NewUseCases(repo)

	tokenManager, err := auth.NewJWTManager("supersecret", 24*time.Hour)
	if err != nil {
		t.Fatalf("failed to initialize token manager: %v", err)
	}

	token, err := tokenManager.NewToken(1)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	handler := handlers.NewHandler(
		useCases.User,
		useCases.Transaction,
		nil,
		nil,
		tokenManager,
	)

	r := router.NewRouter(handler, tokenManager)

	reqBody := map[string]interface{}{
		"ToUser": "bob",
		"Amount": 100,
	}

	reqBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(reqBytes))

	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var senderBalance int

	err = testDB.QueryRow("SELECT balance FROM users WHERE id = $1", 1).Scan(&senderBalance)
	assert.NoError(t, err)
	assert.Equal(t, 900, senderBalance)

	var receiverBalance int

	err = testDB.QueryRow("SELECT balance FROM users WHERE id = $1", 2).Scan(&receiverBalance)
	assert.NoError(t, err)
	assert.Equal(t, 600, receiverBalance)

	var transactionCount int

	query := `
        SELECT COUNT(*) 
        FROM transactions 
        WHERE sender_id = $1 AND receiver_id = $2 AND amount = $3
    `
	err = testDB.QueryRow(query, 1, 2, 100).Scan(&transactionCount)
	assert.NoError(t, err)
	assert.Equal(t, 1, transactionCount)
}
