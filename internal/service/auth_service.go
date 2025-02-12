package service

import (
	"errors"
	"time"

	"github.com/KazikovAP/merch_store/internal/model"
	"github.com/KazikovAP/merch_store/internal/repository"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
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
	user, err := a.userRepo.GetByUsername(username)
	if err != nil {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		newUser := &model.User{
			Username: username,
			Password: string(hashedPassword),
			Coins:    1000,
		}

		if err := a.userRepo.Create(newUser); err != nil {
			return "", err
		}

		user = newUser
	} else {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return "", errors.New("invalid credentials")
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(a.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
