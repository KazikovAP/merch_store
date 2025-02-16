package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/KazikovAP/merch_store/internal/api/http/models"
)

// /api/auth.
func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	var req models.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request format")
		return
	}

	// Получаем или создаем пользователя
	user, err := h.userUseCase.GetByUsername(r.Context(), req.Username)
	if err != nil {
		// Если пользователь не найден - создаем нового
		user, err = h.userUseCase.Register(r.Context(), req.Username)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to create user")
			return
		}
	}

	token, err := h.tokenManager.NewToken(user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create token")
		return
	}

	writeJSON(w, http.StatusOK, models.AuthResponse{Token: token})
}
