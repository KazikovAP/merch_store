package router

import (
	"github.com/KazikovAP/merch_store/internal/config"
	"github.com/KazikovAP/merch_store/internal/handlers"
	"github.com/KazikovAP/merch_store/internal/middleware"
	"github.com/gorilla/mux"
)

func SetupRouter(h *handlers.Handler, _ config.ServerConfig, authCfg config.AuthConfig) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/auth", h.Auth).Methods("POST")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.JwtMiddleware(authCfg.JWTSecret))
	api.HandleFunc("/info", h.Info).Methods("GET")
	api.HandleFunc("/sendCoin", h.SendCoin).Methods("POST")
	api.HandleFunc("/buy/{item}", h.Buy).Methods("GET")

	return r
}
