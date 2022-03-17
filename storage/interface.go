package storage

import (
	"context"

	"github.com/sbxb/loyalty/models"
)

// Storage provides API for writing/reading URLs to/from a data store
type Storage interface {
	AddUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, user *models.User) (*models.User, error)
	AddOrder(ctx context.Context, order *models.Order, userID int) error
	GetOrders(ctx context.Context, userID int) ([]*models.Order, error)
	GetBalance(ctx context.Context, userID int) (models.Balance, error)
	GetWithdrawals(ctx context.Context, userID int) ([]*models.WithdrawalInfo, error)
	GetUnprocessedOrders(ctx context.Context, limit int) ([]*models.Order, error)
	UpdateOrderStatus(ctx context.Context, ar *models.AccrualResponse) error
	ProcessOrder(ctx context.Context, ar *models.AccrualResponse) error
	ProcessWithdraw(ctx context.Context, wr *models.WithdrawRequest, userID int) error
	Close() error
}
