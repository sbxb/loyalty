package order

import (
	"context"
	"testing"

	"github.com/sbxb/loyalty/models"
	"github.com/sbxb/loyalty/storage/inmemory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterOrder(t *testing.T) {
	store, _ := inmemory.NewMapStorage() // NewMapStorage() never returns non-nil error

	orderService := NewOrderService(store)

	tests := []struct {
		name      string
		order     models.Order
		userID    int
		wantError bool
		wantCode  int
	}{
		{
			name:      "First order, first user",
			order:     models.Order{Number: "12345678903"},
			userID:    1,
			wantError: false,
			wantCode:  202,
		},
		{
			name:      "Same order, same user",
			order:     models.Order{Number: "12345678903"},
			userID:    1,
			wantError: true,
			wantCode:  200,
		},
		{
			name:      "Same order, another user",
			order:     models.Order{Number: "12345678903"},
			userID:    2,
			wantError: true,
			wantCode:  409,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := orderService.RegisterOrder(context.Background(), &tt.order, tt.userID)
			if tt.wantError {
				var orderError *OrderError
				require.ErrorAs(t, err, &orderError)
				assert.Equal(t, tt.wantCode, orderError.Code)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestListOrders(t *testing.T) {
	store, _ := inmemory.NewMapStorage() // NewMapStorage() never returns non-nil error

	orderService := NewOrderService(store)

	// Load some orders
	initData := []struct {
		order  models.Order
		userID int
	}{
		{
			order:  models.Order{Number: "12345678903"},
			userID: 1,
		},
		{
			order:  models.Order{Number: "49927398716"},
			userID: 1,
		},
	}
	for _, entry := range initData {
		_ = orderService.RegisterOrder(context.Background(), &entry.order, entry.userID)
	}

	// Test cases
	tests := []struct {
		name           string
		userID         int
		wantError      bool
		wantListLength int
		wantCode       int
	}{
		{
			name:           "First user",
			userID:         1,
			wantError:      false,
			wantListLength: 2,
			wantCode:       0,
		},
		{
			name:           "Second user",
			userID:         2,
			wantError:      true,
			wantListLength: 0,
			wantCode:       204,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orders, err := orderService.ListOrders(context.Background(), tt.userID)
			if tt.wantError {
				var orderError *OrderError
				require.ErrorAs(t, err, &orderError)
				assert.Equal(t, tt.wantCode, orderError.Code)
			} else {
				assert.Equal(t, tt.wantListLength, len(orders))
			}
		})
	}
}
