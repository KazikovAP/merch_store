package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/KazikovAP/merch_store/internal/api/http/middleware"
	"github.com/KazikovAP/merch_store/internal/api/http/models"
)

// /api/sendCoin.
func (h *Handler) SendCoin(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		writeError(w, http.StatusUnauthorized, "user id not found in context")
		return
	}

	var req models.SendCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request format")
		return
	}

	// Получаем получателя
	receiver, err := h.userUseCase.GetByUsername(r.Context(), req.ToUser)
	if err != nil {
		writeError(w, http.StatusBadRequest, "receiver not found")
		return
	}

	// Выполняем перевод
	err = h.transactionUseCase.Transfer(r.Context(), userID, receiver.ID, req.Amount)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
