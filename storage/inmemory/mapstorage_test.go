package inmemory_test

import (
	"context"
	"testing"

	"github.com/sbxb/loyalty/models"
	"github.com/sbxb/loyalty/storage"
	"github.com/sbxb/loyalty/storage/inmemory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddGetUser(t *testing.T) {
	user := &models.User{
		Login: "user",
		Hash:  "abcdef",
		ID:    1,
	}

	store, _ := inmemory.NewMapStorage() // NewMapStorage never returns non-nil error

	err := store.AddUser(context.Background(), user)
	require.NoError(t, err)

	userReturned, err := store.GetUser(context.Background(), user)
	require.NoError(t, err)

	assert.Equal(t, userReturned, user)

	//store.DumpUser()
}

func TestAddUserTwice(t *testing.T) {
	user := &models.User{
		Login: "user",
		Hash:  "abcdef",
	}

	store, _ := inmemory.NewMapStorage() // NewMapStorage never returns non-nil error

	err := store.AddUser(context.Background(), user)
	require.NoError(t, err)

	err = store.AddUser(context.Background(), user)
	require.ErrorIs(t, err, storage.ErrLoginAlreadyExists)
}

func TestGetNonexistentUser(t *testing.T) {
	user := &models.User{
		Login: "usernonexistent",
		Hash:  "nonexistent",
	}

	store, _ := inmemory.NewMapStorage() // NewMapStorage never returns non-nil error

	userReturned, err := store.GetUser(context.Background(), user)
	assert.Nil(t, userReturned)
	require.ErrorIs(t, err, storage.ErrLoginMissing)
}

func TestAddOrder(t *testing.T) {
	order := &models.Order{
		Number: "12345678903",
		Status: models.OrderStatusNew,
	}
	userID := 1

	store, _ := inmemory.NewMapStorage() // NewMapStorage never returns non-nil error

	err := store.AddOrder(context.Background(), order, userID)
	require.NoError(t, err)
}

func TestAddOrderExistsOwn(t *testing.T) {
	order := &models.Order{
		Number: "12345678903",
		Status: models.OrderStatusNew,
	}
	userID := 1

	store, _ := inmemory.NewMapStorage() // NewMapStorage never returns non-nil error

	err := store.AddOrder(context.Background(), order, userID)
	require.NoError(t, err)

	err = store.AddOrder(context.Background(), order, userID)
	var orderExistsError *storage.ExistingOrderError
	require.ErrorAs(t, err, &orderExistsError)
	orderExistsError = err.(*storage.ExistingOrderError)
	assert.Equal(t, orderExistsError.UserID, userID)
}

func TestAddOrderExistsNotOwn(t *testing.T) {
	order := &models.Order{
		Number: "12345678903",
		Status: models.OrderStatusNew,
	}
	userID := 1

	store, _ := inmemory.NewMapStorage() // NewMapStorage never returns non-nil error

	err := store.AddOrder(context.Background(), order, userID)
	require.NoError(t, err)

	err = store.AddOrder(context.Background(), order, 2)
	var orderExistsError *storage.ExistingOrderError
	require.ErrorAs(t, err, &orderExistsError)
	orderExistsError = err.(*storage.ExistingOrderError)
	assert.NotEqual(t, orderExistsError.UserID, 2)

	//store.DumpUser()
	//store.DumpOrder()
}

func TestGetOrdersExistent(t *testing.T) {
	orders := []*models.Order{
		{
			Number:  "12345",
			Status:  "NEW",
			Accrual: 0,
		},
		{
			Number:  "67890",
			Status:  "NEW",
			Accrual: 0,
		},
	}

	userID := 1

	store, _ := inmemory.NewMapStorage() // NewMapStorage never returns non-nil error

	for _, order := range orders {
		err := store.AddOrder(context.Background(), order, userID)
		require.NoError(t, err)
	}

	ordersReturned, err := store.GetOrders(context.Background(), userID)
	require.NoError(t, err)
	assert.NotEmpty(t, ordersReturned)
	for i := range ordersReturned {
		assert.True(t, ordersEqual(ordersReturned[i], orders[i]))
	}

}

func TestGetOrdersNonExistent(t *testing.T) {
	userID := 1
	store, _ := inmemory.NewMapStorage() // NewMapStorage never returns non-nil error

	ordersReturned, err := store.GetOrders(context.Background(), userID)
	require.NoError(t, err)
	assert.Empty(t, ordersReturned)
}

func ordersEqual(first, second *models.Order) bool {
	return first.Number == second.Number &&
		first.Status == second.Status &&
		first.Accrual == second.Accrual
}
