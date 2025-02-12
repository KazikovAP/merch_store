package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/KazikovAP/merch_store/internal/middleware"
	"github.com/KazikovAP/merch_store/internal/model"
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
	var req model.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.authService.Authenticate(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	resp := model.AuthResponse{Token: token}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// /api/info.
func (h *Handler) Info(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(middleware.UsernameKey).(string)

	info, err := h.userService.GetUserInfo(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// /api/sendCoin.
func (h *Handler) SendCoin(w http.ResponseWriter, r *http.Request) {
	var req model.SendCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fromUsername := r.Context().Value(middleware.UsernameKey).(string)
	if err := h.coinService.TransferCoins(fromUsername, req.ToUser, req.Amount); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte("success")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// /api/buy/{item}.
func (h *Handler) Buy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	item := vars["item"]

	username := r.Context().Value(middleware.UsernameKey).(string)
	if err := h.purchaseService.PurchaseItem(username, item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte("success")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
