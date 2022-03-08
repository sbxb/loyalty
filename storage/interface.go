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
	Close() error
}
