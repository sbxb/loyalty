package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/sbxb/loyalty/config"
	"github.com/sbxb/loyalty/internal/logger"
	"github.com/sbxb/loyalty/models"
	"github.com/sbxb/loyalty/services/accrual"
	"github.com/sbxb/loyalty/services/auth"
	"github.com/sbxb/loyalty/services/order"
	"github.com/sbxb/loyalty/storage"
)

// URLHandler defines a container for handlers and their dependencies
type URLHandler struct {
	store   storage.Storage
	config  config.Config
	auth    *auth.AuthService
	ord     *order.OrderService
	accrual *accrual.SimpleAccrualService
}

func NewURLHandler(st storage.Storage, cfg config.Config) URLHandler {
	return URLHandler{
		store:   st,
		config:  cfg,
		auth:    auth.NewAuthService(st),
		ord:     order.NewOrderService(st),
		accrual: accrual.NewSimpleAccrualService(st, cfg.AccrualAddress),
	}
}

// UserRegister process POST /api/user/register request
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

	// http.StatusOK sent implicitly
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

	// http.StatusOK sent implicitly
}

// UserPostOrder process POST /api/user/orders request
func (uh URLHandler) UserPostOrder(w http.ResponseWriter, r *http.Request) {
	logger.Info("UserPostOrder hit by POST /api/user/orders")
	order, orderErr := ReadOrderNumberFromBody(r.Body)
	if orderErr != nil {
		http.Error(w, orderErr.Error(), orderErr.Code)
		return
	}

	userID := auth.GetUserID(r.Context())
	if orderRegErr := uh.ord.RegisterOrder(r.Context(), order, userID); orderRegErr != nil {
		http.Error(w, orderRegErr.Error(), orderRegErr.Code)
		return
	}

	go func() {
		retry := uh.accrual.DoAccrualStuff(order.Number)
		if retry {
			time.Sleep(500 * time.Microsecond)
			_ = uh.accrual.DoAccrualStuff(order.Number)
		}
	}()

	w.WriteHeader(http.StatusAccepted)
}

// UserGetOrders process GET /api/user/orders request
func (uh URLHandler) UserGetOrders(w http.ResponseWriter, r *http.Request) {
	logger.Info("UserGetOrders hit by GET /api/user/orders")
	userID := auth.GetUserID(r.Context())

	orderList, orderRegErr := uh.ord.ListOrders(r.Context(), userID)
	if orderRegErr != nil {
		http.Error(w, orderRegErr.Error(), orderRegErr.Code)
		return
	}

	jr, err := json.Marshal(orderList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jr)
}

// UserGetBalance process GET /api/user/balance request
func (uh URLHandler) UserGetBalance(w http.ResponseWriter, r *http.Request) {
	logger.Info("UserGetBalance hit by GET /api/user/balance")
	userID := auth.GetUserID(r.Context())

	balance, err := uh.store.GetBalance(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jr, err := json.Marshal(balance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jr)
}

// UserBalanceWithdraw process POST /api/user/balance/withdraw request
func (uh URLHandler) UserBalanceWithdraw(w http.ResponseWriter, r *http.Request) {
	logger.Info("UserBalanceWithdraw hit by POST /api/user/balance/withdraw")
	userID := auth.GetUserID(r.Context())

	req, err := models.ReadWithdrawRequestFromBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !req.Validate() {
		http.Error(w, "wrong number format", http.StatusUnprocessableEntity)
		return
	}

	err = uh.store.ProcessWithdraw(r.Context(), req, userID)
	if err != nil {
		if errors.Is(err, storage.ErrInsufficientFunds) {
			http.Error(w, err.Error(), http.StatusPaymentRequired)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// http.StatusOK sent implicitly
}

// UserGetWithdrawals process GET /api/user/balance/withdrawals request
func (uh URLHandler) UserGetWithdrawals(w http.ResponseWriter, r *http.Request) {
	logger.Info("UserGetWithdrawals hit by GET /api/user/balance/withdrawals")
	userID := auth.GetUserID(r.Context())

	withdrawals, err := uh.store.GetWithdrawals(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(withdrawals) == 0 {
		http.Error(w, "no withdrawals", http.StatusNoContent)
		return
	}

	jr, err := json.Marshal(withdrawals)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jr)
}
