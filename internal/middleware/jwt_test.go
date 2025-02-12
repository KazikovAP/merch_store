package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/KazikovAP/merch_store/internal/middleware"
	"github.com/golang-jwt/jwt"
)

func createTestToken(secret, username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
	})

	tokenString, _ := token.SignedString([]byte(secret))

	return tokenString
}

func TestJwtMiddleware_MissingAuthHeader(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if _, err := w.Write([]byte("ok")); err != nil {
			t.Errorf("error writing response: %v", err)
		}
	})

	mw := middleware.JwtMiddleware("testsecret")
	handlerToTest := mw(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rr := httptest.NewRecorder()

	handlerToTest.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "Missing Authorization header") {
		t.Errorf("expected error message about missing header, got %s", rr.Body.String())
	}
}

func TestJwtMiddleware_InvalidToken(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if _, err := w.Write([]byte("ok")); err != nil {
			t.Errorf("error writing response: %v", err)
		}
	})

	mw := middleware.JwtMiddleware("testsecret")
	handlerToTest := mw(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	req.Header.Set("Authorization", "Bearer invalidtoken")

	rr := httptest.NewRecorder()

	handlerToTest.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "Invalid token") {
		t.Errorf("expected error message about invalid token, got %s", rr.Body.String())
	}
}

func TestJwtMiddleware_ValidToken(t *testing.T) {
	secret := "testsecret"
	username := "testuser"
	token := createTestToken(secret, username)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userFromCtx, ok := r.Context().Value(middleware.UsernameKey).(string)
		if !ok {
			http.Error(w, "username not found in context", http.StatusInternalServerError)
			return
		}

		if _, err := w.Write([]byte(userFromCtx)); err != nil {
			t.Errorf("error writing response: %v", err)
		}
	})

	mw := middleware.JwtMiddleware(secret)
	handlerToTest := mw(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()

	handlerToTest.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if !strings.Contains(rr.Body.String(), username) {
		t.Errorf("expected response to contain username %s, got %s", username, rr.Body.String())
	}
}
