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

	user    map[string]string // login -> id|login|hash
	order   map[string]string // number -> status|accrual|uploaded_at|user_id
	balance map[int]string    // user_id -> current|withdrawn
}

// MapStorage implements Storage interface
var _ storage.Storage = (*MapStorage)(nil)

func NewMapStorage() (*MapStorage, error) {
	user := make(map[string]string)
	order := make(map[string]string)
	balance := make(map[int]string)
	return &MapStorage{user: user, order: order, balance: balance}, nil
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
	uid := len(ms.user) + 1
	ms.user[user.Login] = fmt.Sprintf("%d|%s|%s", uid, user.Login, user.Hash)
	ms.balance[uid] = fmt.Sprintf("%d|%d", 0, 0)

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

		var acc int64
		if acc, err = strconv.ParseInt(parts[1], 10, 64); err != nil {
			return nil, fmt.Errorf("MapStorage: GetOrders: %v", err)
		}
		order.Accrual = models.Money(acc)
		if order.UploadedAt, err = time.Parse(time.RFC3339, parts[2]); err != nil {
			return nil, fmt.Errorf("MapStorage: GetOrders: %v", err)
		}
		res = append(res, order)
	}

	return res, nil
}

func (ms *MapStorage) GetBalance(ctx context.Context, userID int) (models.Balance, error) {
	ms.Lock()
	defer ms.Unlock()

	balance := models.Balance{}

	for uid, payload := range ms.balance {
		if uid == userID {
			var current, withdrawn int64
			parts := strings.SplitN(payload, "|", 2)
			current, _ = strconv.ParseInt(parts[0], 10, 64)
			balance.Current = models.Money(current)
			withdrawn, _ = strconv.ParseInt(parts[1], 10, 64)
			balance.Withdrawn = models.Money(withdrawn)
		}
	}

	return balance, nil
}

func (ms *MapStorage) GetWithdrawals(ctx context.Context, userID int) ([]*models.WithdrawalInfo, error) {
	ms.Lock()
	defer ms.Unlock()
	// TODO implementation
	res := []*models.WithdrawalInfo{}

	return res, nil
}

func (ms *MapStorage) GetUnprocessedOrders(ctx context.Context, limit int) ([]*models.Order, error) {
	ms.Lock()
	defer ms.Unlock()

	res := []*models.Order{}

	// TODO Sort orders by upload time ?
	for key, payload := range ms.order {
		parts := strings.SplitN(payload, "|", 4)
		status := parts[0]
		if status != models.OrderStatusNew && status != models.OrderStatusProcessing {
			continue
		}

		order := &models.Order{}
		order.Number = key // All we are interested in is the current order's number

		res = append(res, order)
		if len(res) >= limit {
			break
		}
	}

	return res, nil
}

func (ms *MapStorage) UpdateOrderStatus(ctx context.Context, ar *models.AccrualResponse) error {
	//

	return nil
}

func (ms *MapStorage) ProcessOrder(ctx context.Context, ar *models.AccrualResponse) error {
	//

	return nil
}

func (ms *MapStorage) ProcessWithdraw(ctx context.Context, wr *models.WithdrawRequest, userID int) error {
	//

	return nil
}

func (ms *MapStorage) Close() error {
	return nil
}

func (ms *MapStorage) DumpUser() {
	for key, payload := range ms.user {
		fmt.Println(key, "=>", payload)
	}
}

func (ms *MapStorage) DumpOrder() {
	for key, payload := range ms.order {
		fmt.Println(key, "=>", payload)
	}
}

func (ms *MapStorage) DumpBalance() {
	for key, payload := range ms.balance {
		fmt.Println(key, "=>", payload)
	}
}
