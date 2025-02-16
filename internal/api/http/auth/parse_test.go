package auth_test

import (
	"testing"
	"time"

	"github.com/KazikovAP/merch_store/internal/api/http/auth"
	"github.com/stretchr/testify/assert"
)

const signingKey = "secret_key"

// Тест на успешный парсинг токена.
func TestJWTManager_Parse_Success(t *testing.T) {
	tokenTTL := 1 * time.Hour

	manager, err := auth.NewJWTManager(signingKey, tokenTTL)
	assert.NoError(t, err)

	userID := 1
	token, err := manager.NewToken(userID)
	assert.NoError(t, err)

	parsedUserID, err := manager.Parse(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, parsedUserID)
}

// Тест на ошибку при парсинге невалидного токена.
func TestJWTManager_Parse_InvalidToken_Error(t *testing.T) {
	tokenTTL := 1 * time.Hour

	manager, err := auth.NewJWTManager(signingKey, tokenTTL)
	assert.NoError(t, err)

	invalidToken := "invalid_token"
	parsedUserID, err := manager.Parse(invalidToken)
	assert.Error(t, err)
	assert.Equal(t, 0, parsedUserID)
}

// Тест на ошибку при парсинге токена с неверным signing key.
func TestJWTManager_Parse_WrongSigningKey_Error(t *testing.T) {
	tokenTTL := 1 * time.Hour

	manager, err := auth.NewJWTManager(signingKey, tokenTTL)
	assert.NoError(t, err)

	userID := 1
	token, err := manager.NewToken(userID)
	assert.NoError(t, err)

	wrongSigningKey := "wrong_secret_key"
	wrongManager, err := auth.NewJWTManager(wrongSigningKey, tokenTTL)
	assert.NoError(t, err)

	parsedUserID, err := wrongManager.Parse(token)
	assert.Error(t, err)
	assert.Equal(t, 0, parsedUserID)
}

// Тест на ошибку при парсинге просроченного токена.
func TestJWTManager_Parse_ExpiredToken_Error(t *testing.T) {
	tokenTTL := -1 * time.Hour

	manager, err := auth.NewJWTManager(signingKey, tokenTTL)
	assert.NoError(t, err)

	userID := 1
	token, err := manager.NewToken(userID)
	assert.NoError(t, err)

	parsedUserID, err := manager.Parse(token)
	assert.Error(t, err)
	assert.Equal(t, 0, parsedUserID)
}
