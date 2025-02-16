package auth

import "github.com/golang-jwt/jwt"

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

type TokenManager interface {
	NewToken(userID int) (string, error)
	Parse(accessToken string) (int, error)
}
