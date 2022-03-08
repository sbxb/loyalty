package handlers

import (
	"io"
	"net/http"

	"github.com/sbxb/loyalty/config"
	"github.com/sbxb/loyalty/internal/logger"
	"github.com/sbxb/loyalty/models"
	"github.com/sbxb/loyalty/services/auth"
	"github.com/sbxb/loyalty/storage"
)

// URLHandler defines a container for handlers and their dependencies
type URLHandler struct {
	store  storage.Storage
	config config.Config
	auth   *auth.AuthService
}

func NewURLHandler(st storage.Storage, cfg config.Config) URLHandler {
	return URLHandler{
		store:  st,
		config: cfg,
		auth:   auth.NewAuthService(st),
	}
}

// UserRegister process POST /api/user/register request
// ... Регистрация производится по паре логин/пароль. Каждый логин должен быть
// уникальным. После успешной регистрации должна происходить автоматическая
// аутентификация пользователя ...
func (uh URLHandler) UserRegister(w http.ResponseWriter, r *http.Request) {
	logger.Info("UserRegister hit by POST /api/user/register")

	user, err := models.ReadUserFromBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authErr := uh.auth.RegisterUser(r.Context(), user)
	if authErr != nil {
		http.Error(w, authErr.Error(), authErr.Code)
		return
	}

	authUser, authErr := uh.auth.LoginUser(r.Context(), user)
	if authErr != nil {
		http.Error(w, authErr.Error(), authErr.Code)
		return
	}

	if err = uh.auth.SetCookie(w, authUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// send http.StatusOK by default
}

// UserLogin process POST /api/user/login request
func (uh URLHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	logger.Info("UserLogin hit by POST /api/user/login")

	user, err := models.ReadUserFromBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authUser, authErr := uh.auth.LoginUser(r.Context(), user)
	if authErr != nil {
		http.Error(w, authErr.Error(), authErr.Code)
		return
	}

	if err = uh.auth.SetCookie(w, authUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// send http.StatusOK by default
}

// UserPostOrder process POST /api/user/orders request
func (uh URLHandler) UserPostOrder(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Server failed to read the request's body", http.StatusInternalServerError)
		return
	}

	order := models.Order{Number: string(body)}
	if !order.Validate() {
		//
	}

}

// UserGetOrders process GET /api/user/orders request
func (uh URLHandler) UserGetOrders(w http.ResponseWriter, r *http.Request) {
	//
}

// UserGetBalance process GET /api/user/balance request
func (uh URLHandler) UserGetBalance(w http.ResponseWriter, r *http.Request) {
	//
}

// UserBalanceWithdraw process POST /api/user/balance/withdraw request
func (uh URLHandler) UserBalanceWithdraw(w http.ResponseWriter, r *http.Request) {
	//
}

// UserGetWithdrawals process GET /api/user/balance/withdrawals request
func (uh URLHandler) UserGetWithdrawals(w http.ResponseWriter, r *http.Request) {
	//
}
