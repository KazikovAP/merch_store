package dto

type AuthRequest struct {
	Username string
	Password string
}

type AuthResponse struct {
	Token string
}
