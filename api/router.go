package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbxb/loyalty/internal/logger"
)

func NewRouter( /*store storage.Storage, */ /*cfg config.Config*/ ) http.Handler {
	router := chi.NewRouter()
	logger.Info("Router created")
	return router
}
