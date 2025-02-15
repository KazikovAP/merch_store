package service

import (
	"errors"
	"time"

	"github.com/KazikovAP/merch_store/internal/model/domain"
	"github.com/KazikovAP/merch_store/internal/repository"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	defaultCoins     = 1000
	tokenExpiryHours = 72
)

type AuthService interface {
	Authenticate(username, password string) (string, error)
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (a *authService) Authenticate(username, password string) (string, error) {
	if username == "" || password == "" {
		return "", errors.New("username and password must not be empty")
	}

	user, err := a.userRepo.GetByUsername(username)
	if err != nil {
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if hashErr != nil {
			return "", errors.New("failed to hash password")
		}

		newUser := &domain.User{
			Username: username,
			Password: string(hashedPassword),
			Coins:    defaultCoins,
		}

		if createErr := a.userRepo.Create(newUser); createErr != nil {
			return "", createErr
		}

		user = newUser
	} else {
		if compareErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); compareErr != nil {
			return "", errors.New("invalid username or password")
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * tokenExpiryHours).Unix(),
	})

	tokenString, err := token.SignedString([]byte(a.jwtSecret))
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return tokenString, nil
}
