package auth_test

import (
	"testing"
	"time"

	"github.com/KazikovAP/merch_store/internal/api/http/auth"
	"github.com/stretchr/testify/assert"
)

// Тест на успешное создание токена.
func TestJWTManager_NewToken_Success(t *testing.T) {
	signingKey := "secret_key"
	tokenTTL := 1 * time.Hour

	manager, err := auth.NewJWTManager(signingKey, tokenTTL)
	assert.NoError(t, err)

	userID := 1
	token, err := manager.NewToken(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

// Тест на ошибку при создании токена с пустым signing key.
func TestJWTManager_NewToken_EmptySigningKey_Error(t *testing.T) {
	signingKey := ""
	tokenTTL := 1 * time.Hour

	manager, err := auth.NewJWTManager(signingKey, tokenTTL)
	assert.Error(t, err)
	assert.Nil(t, manager)
}
