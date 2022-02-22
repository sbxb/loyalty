package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbxb/loyalty/api/handlers"
	"github.com/sbxb/loyalty/config"
	"github.com/sbxb/loyalty/internal/logger"
)

func NewRouter( /*store storage.Storage, */ cfg config.Config) http.Handler {
	router := chi.NewRouter()
	logger.Info("Router created")
	urlHandler := handlers.NewURLHandler( /*store, */ cfg)

	router.Post("/api/user/register", urlHandler.UserRegister)
	router.Post("/api/user/login", urlHandler.UserLogin)

	return router
}
