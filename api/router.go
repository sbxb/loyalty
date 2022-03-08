package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbxb/loyalty/api/handlers"
	"github.com/sbxb/loyalty/config"
	"github.com/sbxb/loyalty/internal/logger"
	"github.com/sbxb/loyalty/storage"
)

func NewRouter(store storage.Storage, cfg config.Config) http.Handler {
	router := chi.NewRouter()
	logger.Info("Router created")

	urlHandler := handlers.NewURLHandler(store, cfg)

	router.Post("/api/user/register", urlHandler.UserRegister)
	router.Post("/api/user/login", urlHandler.UserLogin)

	router.Post("/api/user/orders", urlHandler.UserPostOrder)
	router.Get("/api/user/orders", urlHandler.UserGetOrders)

	router.Get("/api/user/balance", urlHandler.UserGetBalance)
	router.Post("/api/user/balance/withdraw", urlHandler.UserBalanceWithdraw)
	router.Get("/api/user/balance/withdrawals", urlHandler.UserGetWithdrawals)

	return router
}
