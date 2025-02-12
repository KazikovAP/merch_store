package service_test

import (
	"testing"
	"time"

	"github.com/KazikovAP/merch_store/internal/model"
	"github.com/KazikovAP/merch_store/internal/service"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthenticate_NewUser(t *testing.T) {
	repo := &mockUserRepo{user: nil}
	jwtSecret := "testsecret1"
	authService := service.NewAuthService(repo, jwtSecret)

	username := "newuser"
	password := "password123"

	token, err := authService.Authenticate(username, password)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token == "" {
		t.Fatalf("expected token, got empty string")
	}

	user, err := repo.GetByUsername(username)
	if err != nil {
		t.Fatalf("expected user to be created, got error: %v", err)
	}

	if user.Coins != 1000 {
		t.Errorf("expected initial coins 1000, got %d", user.Coins)
	}

	parsedToken, err := jwt.Parse(token, func(_ *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !parsedToken.Valid {
		t.Fatalf("expected valid token, got error: %v", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("expected MapClaims, got different type")
	}

	if claims["username"] != username {
		t.Errorf("expected username %s in token, got %v", username, claims["username"])
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		t.Fatalf("expected exp claim to be a number")
	}

	expTime := time.Unix(int64(expFloat), 0)
	if time.Until(expTime) < 70*time.Hour/3 {
		t.Errorf("expected token expiration around 72 hours, got %v left", time.Until(expTime))
	}
}

func TestAuthenticate_ValidCredentials(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	repo := &mockUserRepo{user: &model.User{
		ID:       1,
		Username: "existinguser",
		Password: string(hashedPassword),
		Coins:    1000,
	}}

	jwtSecret := "testsecret2"
	authService := service.NewAuthService(repo, jwtSecret)

	token, err := authService.Authenticate("existinguser", "password123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token == "" {
		t.Fatalf("expected token, got empty string")
	}

	parsedToken, err := jwt.Parse(token, func(_ *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !parsedToken.Valid {
		t.Fatalf("expected valid token, got error: %v", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("expected MapClaims, got different type")
	}

	if claims["username"] != "existinguser" {
		t.Errorf("expected username %s in token, got %v", "existinguser", claims["username"])
	}
}

func TestAuthenticate_InvalidPassword(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	repo := &mockUserRepo{user: &model.User{
		ID:       1,
		Username: "existinguser",
		Password: string(hashedPassword),
		Coins:    1000,
	}}

	jwtSecret := "testsecret3"
	authService := service.NewAuthService(repo, jwtSecret)

	token, err := authService.Authenticate("existinguser", "wrongpassword")
	if err == nil {
		t.Fatalf("expected error due to invalid credentials, got nil")
	}

	if token != "" {
		t.Errorf("expected empty token on failure, got %s", token)
	}
}
