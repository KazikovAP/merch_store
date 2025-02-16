package handlers

import (
	"net/http"

	"github.com/KazikovAP/merch_store/internal/api/http/middleware"
	"github.com/gorilla/mux"
)

// /api/buy/{item}.
func (h *Handler) Buy(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		writeError(w, http.StatusUnauthorized, "user id not found in context")
		return
	}

	vars := mux.Vars(r)
	merchName := vars["item"]

	// Получаем информацию о товаре
	merch, err := h.merchUseCase.GetByName(r.Context(), merchName)
	if err != nil {
		writeError(w, http.StatusBadRequest, "item not found")
		return
	}

	// Выполняем покупку
	err = h.purchaseUseCase.Purchase(r.Context(), userID, 1, merch.Name)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
