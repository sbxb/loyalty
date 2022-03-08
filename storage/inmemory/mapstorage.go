package inmemory

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sbxb/loyalty/models"
	"github.com/sbxb/loyalty/storage"
)

// MapStorage defines a simple in-memory storage implemented as a wrapper
// around Go maps
type MapStorage struct {
	sync.RWMutex

	user  map[string]string // login -> id|login|hash
	order map[string]string // number -> status|accrual|uploaded_at|user_id
}

// MapStorage implements Storage interface
var _ storage.Storage = (*MapStorage)(nil)

func NewMapStorage() (*MapStorage, error) {
	user := make(map[string]string)
	order := make(map[string]string)
	return &MapStorage{user: user, order: order}, nil
}

func (ms *MapStorage) AddUser(ctx context.Context, user *models.User) error {
	ms.Lock()
	defer ms.Unlock()

	// check unique constraint on login
	for key := range ms.user {
		if key == user.Login {
			return storage.ErrLoginAlreadyExists
		}
	}

	// add new user
	ms.user[user.Login] = fmt.Sprintf("%d|%s|%s", len(ms.user)+1, user.Login, user.Hash)

	return nil
}

func (ms *MapStorage) GetUser(ctx context.Context, user *models.User) (*models.User, error) {
	ms.Lock()
	defer ms.Unlock()
	dbUser := &models.User{}

	for key, payload := range ms.user {
		if key == user.Login {
			parts := strings.SplitN(payload, "|", 3)
			dbUser.ID, _ = strconv.Atoi(parts[0])
			dbUser.Login = parts[1]
			dbUser.Hash = parts[2]
			return dbUser, nil
		}
	}
	return nil, storage.ErrLoginMissing
}

func (ms *MapStorage) AddOrder(ctx context.Context, order *models.Order, userID int) error {
	ms.Lock()
	defer ms.Unlock()

	// check unique constraint on number
	for key, payload := range ms.order {
		if key == order.Number {
			parts := strings.SplitN(payload, "|", 4)
			uid, _ := strconv.Atoi(parts[3])
			return storage.NewExistingOrderError(uid)
		}
	}
	currDate := time.Now().Format(time.RFC3339)
	ms.order[order.Number] = fmt.Sprintf("%s|%d|%s|%d", order.Status, order.Accrual, currDate, userID)

	return nil
}

func (ms *MapStorage) GetOrders(ctx context.Context, userID int) ([]*models.Order, error) {
	ms.Lock()
	defer ms.Unlock()

	res := []*models.Order{}

	// TODO Sort orders by upload time ?
	for key, payload := range ms.order {
		parts := strings.SplitN(payload, "|", 4)
		uid, _ := strconv.Atoi(parts[3])
		if uid != userID {
			continue
		}
		var err error
		order := &models.Order{}
		order.Number = key
		order.Status = parts[0]
		if order.Accrual, err = strconv.ParseInt(parts[1], 10, 64); err != nil {
			return nil, fmt.Errorf("MapStorage: GetOrders: %v", err)
		}
		if order.UploadedAt, err = time.Parse(time.RFC3339, parts[2]); err != nil {
			return nil, fmt.Errorf("MapStorage: GetOrders: %v", err)
		}
		res = append(res, order)
	}

	return res, nil
}

func (ms *MapStorage) Close() error {
	return nil
}

func (ms *MapStorage) DumpUser() {
	for _, payload := range ms.user {
		fmt.Println(payload)
	}
}

func (ms *MapStorage) DumpOrder() {
	for _, payload := range ms.order {
		fmt.Println(payload)
	}
}
