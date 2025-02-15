package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/KazikovAP/merch_store/internal/middleware"
	"github.com/KazikovAP/merch_store/internal/model/dto"
	"github.com/KazikovAP/merch_store/internal/service"
	"github.com/gorilla/mux"
)

type Handler struct {
	authService     service.AuthService
	userService     service.UserService
	coinService     service.CoinService
	purchaseService service.PurchaseService
}

func NewHandler(auth service.AuthService, user service.UserService, coin service.CoinService, purchase service.PurchaseService) *Handler {
	return &Handler{
		authService:     auth,
		userService:     user,
		coinService:     coin,
		purchaseService: purchase,
	}
}

// /api/auth.
func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	var req dto.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Authenticate(req.Username, req.Password)
	if err != nil {
		http.Error(w, "authentication failed", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(dto.AuthResponse{Token: token}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// /api/info.
func (h *Handler) Info(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(middleware.UsernameKey).(string)
	if !ok {
		http.Error(w, "failed to retrieve username from context", http.StatusInternalServerError)
		return
	}

	info, err := h.userService.GetUserInfo(username)
	if err != nil {
		http.Error(w, "failed to get user info", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// /api/sendCoin.
func (h *Handler) SendCoin(w http.ResponseWriter, r *http.Request) {
	var req dto.SendCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	fromUsername, ok := r.Context().Value(middleware.UsernameKey).(string)
	if !ok {
		http.Error(w, "failed to retrieve username from context", http.StatusInternalServerError)
		return
	}

	if err := h.coinService.TransferCoins(fromUsername, req.ToUser, req.Amount); err != nil {
		http.Error(w, "failed to transfer coins", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// /api/buy/{item}.
func (h *Handler) Buy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	item, ok := vars["item"]
	if !ok {
		http.Error(w, "item not specified", http.StatusBadRequest)
		return
	}

	username, ok := r.Context().Value(middleware.UsernameKey).(string)
	if !ok {
		http.Error(w, "failed to retrieve username from context", http.StatusInternalServerError)
		return
	}

	if err := h.purchaseService.PurchaseItem(username, item); err != nil {
		if err.Error() == "insufficient funds" {
			http.Error(w, "insufficient funds", http.StatusPaymentRequired)
			return
		}

		http.Error(w, "purchase failed", http.StatusBadRequest)

		return
	}

	w.WriteHeader(http.StatusOK)
}
