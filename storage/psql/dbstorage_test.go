package psql_test

import (
	"context"
	"testing"

	"github.com/sbxb/loyalty/models"
	"github.com/sbxb/loyalty/storage"
	"github.com/sbxb/loyalty/storage/psql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const dsn = "postgres://loyalty:loyalty@localhost/loyaltytest"

func TestAddGetUser(t *testing.T) {
	user := &models.User{
		Login: "user",
		Hash:  "abcdef",
		ID:    1,
	}

	store, err := psql.NewDBStorage(dsn)
	require.NoError(t, err)
	err = store.TruncateTables()
	require.NoError(t, err)

	err = store.AddUser(context.Background(), user)
	require.NoError(t, err)

	userReturned, err := store.GetUser(context.Background(), user)
	require.NoError(t, err)

	assert.Equal(t, userReturned, user)
}

func TestAddUserTwice(t *testing.T) {
	user := &models.User{
		Login: "user",
		Hash:  "abcdef",
	}

	store, err := psql.NewDBStorage(dsn)
	require.NoError(t, err)
	err = store.TruncateTables()
	require.NoError(t, err)

	err = store.AddUser(context.Background(), user)
	require.NoError(t, err)

	err = store.AddUser(context.Background(), user)
	require.ErrorIs(t, err, storage.ErrLoginAlreadyExists)
}

func TestGetNonexistentUser(t *testing.T) {
	user := &models.User{
		Login: "usernonexistent",
		Hash:  "nonexistent",
	}

	store, err := psql.NewDBStorage(dsn)
	require.NoError(t, err)
	err = store.TruncateTables()
	require.NoError(t, err)

	userReturned, err := store.GetUser(context.Background(), user)
	assert.Nil(t, userReturned)
	require.ErrorIs(t, err, storage.ErrLoginMissing)
}

func TestAddOrder(t *testing.T) {
	order := &models.Order{
		Number: "12345678903",
		Status: "NEW",
	}
	userID := 1

	store, err := psql.NewDBStorage(dsn)
	require.NoError(t, err)
	err = store.TruncateTables()
	require.NoError(t, err)

	user := &models.User{
		Login: "user",
		Hash:  "abcdef",
	}
	err = store.AddUser(context.Background(), user)
	require.NoError(t, err)

	err = store.AddOrder(context.Background(), order, userID)
	require.NoError(t, err)
}

func TestAddOrderExistsOwn(t *testing.T) {
	order := &models.Order{
		Number: "12345678903",
		Status: "NEW",
	}
	userID := 1

	store, err := psql.NewDBStorage(dsn)
	require.NoError(t, err)
	err = store.TruncateTables()
	require.NoError(t, err)

	user := &models.User{
		Login: "user",
		Hash:  "abcdef",
	}
	err = store.AddUser(context.Background(), user)
	require.NoError(t, err)

	err = store.AddOrder(context.Background(), order, userID)
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
		Status: "NEW",
	}
	userID := 1

	store, err := psql.NewDBStorage(dsn)
	require.NoError(t, err)
	err = store.TruncateTables()
	require.NoError(t, err)

	user := &models.User{
		Login: "user",
		Hash:  "abcdef",
	}
	err = store.AddUser(context.Background(), user)
	require.NoError(t, err)

	err = store.AddOrder(context.Background(), order, userID)
	require.NoError(t, err)

	err = store.AddOrder(context.Background(), order, 2)
	var orderExistsError *storage.ExistingOrderError
	require.ErrorAs(t, err, &orderExistsError)
	orderExistsError = err.(*storage.ExistingOrderError)
	assert.NotEqual(t, orderExistsError.UserID, 2)
}
