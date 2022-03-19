package order

import (
	"context"
	"errors"
	"net/http"

	"github.com/sbxb/loyalty/models"
	"github.com/sbxb/loyalty/storage"
)

type OrderService struct {
	store storage.Storage
}

func NewOrderService(st storage.Storage) *OrderService {
	return &OrderService{store: st}
}

type OrderError struct {
	msg  string
	Code int
}

func (oe OrderError) Error() string {
	return oe.msg
}

func NewOrderError(msg string, code int) *OrderError {
	return &OrderError{
		msg:  msg,
		Code: code,
	}
}

func (osv *OrderService) RegisterOrder(ctx context.Context, order *models.Order, userID int) *OrderError {
	order.Status = models.OrderStatusNew
	if err := osv.store.AddOrder(ctx, order, userID); err != nil {
		var eoErr *storage.ExistingOrderError
		if errors.As(err, &eoErr) {
			if eoErr.UserID == userID {
				return NewOrderError("order has already been loaded by the current user", http.StatusOK)
			}
			return NewOrderError("order has already been loaded by another user", http.StatusConflict)
		}
		return NewOrderError(err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (osv *OrderService) ListOrders(ctx context.Context, userID int) ([]*models.Order, *OrderError) {
	orders, err := osv.store.GetOrders(ctx, userID)

	if err != nil {
		return nil, NewOrderError(err.Error(), http.StatusInternalServerError)
	}
	if len(orders) == 0 {
		return nil, NewOrderError("no orders found", http.StatusNoContent)
	}
	// TODO remove accrual if zero
	return orders, nil
}
