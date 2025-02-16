package handlers

import (
	"net/http"

	"github.com/KazikovAP/merch_store/internal/api/http/middleware"
	"github.com/KazikovAP/merch_store/internal/api/http/models"
)

// /api/info.
func (h *Handler) Info(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		writeError(w, http.StatusUnauthorized, "user id not found in context")
		return
	}

	user, err := h.userUseCase.GetByID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get user info")
		return
	}

	purchases, err := h.purchaseUseCase.GetUserPurchases(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get user purchases")
		return
	}

	sentTx, err := h.transactionUseCase.GetSentTransactions(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get sent transactions")
		return
	}

	receivedTx, err := h.transactionUseCase.GetReceivedTransactions(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get received transactions")
		return
	}

	resp := models.InfoResponse{
		Coins:     user.Balance,
		Inventory: mapInventory(purchases),
		CoinHistory: models.CoinHistoryInfo{
			Sent:     mapTransactions(sentTx, false),
			Received: mapTransactions(receivedTx, true),
		},
	}

	writeJSON(w, http.StatusOK, resp)
}
