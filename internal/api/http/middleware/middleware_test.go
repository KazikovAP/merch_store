package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KazikovAP/merch_store/internal/api/http/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var ErrInvalidToken = errors.New("invalid token")

type MockTokenManager struct {
	mock.Mock
}

func (m *MockTokenManager) Parse(accessToken string) (int, error) {
	args := m.Called(accessToken)
	return args.Int(0), args.Error(1)
}

func (m *MockTokenManager) NewToken(userID int) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func TestAuthMiddleware_MissingAuthHeader(t *testing.T) {
	mockTokenManager := &MockTokenManager{}

	middlewareFunc := middleware.AuthMiddleware(mockTokenManager)

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/info", http.NoBody)
	w := httptest.NewRecorder()

	middlewareFunc(handler).ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "empty auth header")
}

func TestAuthMiddleware_InvalidAuthHeaderFormat(t *testing.T) {
	mockTokenManager := &MockTokenManager{}

	middlewareFunc := middleware.AuthMiddleware(mockTokenManager)

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/info", http.NoBody)

	req.Header.Set("Authorization", "InvalidFormat")

	w := httptest.NewRecorder()

	middlewareFunc(handler).ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid auth header")
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	mockTokenManager := &MockTokenManager{}
	mockTokenManager.On("Parse", "invalid_token").Return(0, ErrInvalidToken)

	middlewareFunc := middleware.AuthMiddleware(mockTokenManager)

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/info", http.NoBody)

	req.Header.Set("Authorization", "Bearer invalid_token")

	w := httptest.NewRecorder()

	middlewareFunc(handler).ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid token")
}

func TestAuthMiddleware_Success(t *testing.T) {
	mockTokenManager := &MockTokenManager{}
	mockTokenManager.On("Parse", "valid_token").Return(1, nil)

	middlewareFunc := middleware.AuthMiddleware(mockTokenManager)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.UserIDKey).(int)
		assert.Equal(t, 1, userID)
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/info", http.NoBody)

	req.Header.Set("Authorization", "Bearer valid_token")

	w := httptest.NewRecorder()

	middlewareFunc(handler).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
