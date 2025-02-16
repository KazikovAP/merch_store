package router

import (
	"net/http"

	"github.com/KazikovAP/merch_store/internal/api/http/auth"
	"github.com/KazikovAP/merch_store/internal/api/http/handlers"
	"github.com/KazikovAP/merch_store/internal/api/http/middleware"
	"github.com/gorilla/mux"
)

func NewRouter(h *handlers.Handler, tokenManager auth.TokenManager) *mux.Router {
	r := mux.NewRouter()

	// Public endpoints
	r.HandleFunc("/api/auth", h.Auth).Methods(http.MethodPost)

	// Protected endpoints
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware(tokenManager))

	api.HandleFunc("/info", h.Info).Methods(http.MethodGet)
	api.HandleFunc("/sendCoin", h.SendCoin).Methods(http.MethodPost)
	api.HandleFunc("/buy/{item}", h.Buy).Methods(http.MethodGet)

	return r
}
