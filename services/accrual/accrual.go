package accrual

import (
	"context"
	"sync"
	"time"

	"github.com/sbxb/loyalty/internal/logger"
	"github.com/sbxb/loyalty/models"
	"github.com/sbxb/loyalty/storage"
)

const (
	QueueLength = 500 // Should not be less than 5
	AddTimeout  = 100 * time.Millisecond
)

type Job struct {
	orderNumber string
	status      string
}

type AccrualService struct {
	store          storage.Storage
	address        string
	newJobQueue    chan Job
	updateJobQueue chan Job
	ctx            context.Context // FIXME Wrong, pass as a parameter instead
}

func NewAccrualService(st storage.Storage, address string, ctx context.Context) *AccrualService {
	logger.Info("Accrual Service : created")
	// TODO make context work
	return &AccrualService{
		store:          st,
		newJobQueue:    make(chan Job, QueueLength),
		updateJobQueue: make(chan Job, QueueLength),
		ctx:            ctx,
		address:        address,
	}
}

func (as *AccrualService) AddOrderNumber(orderNumber string) {
	job := Job{
		orderNumber: orderNumber,
		status:      models.OrderStatusNew,
	}

	select {
	case as.newJobQueue <- job:
		logger.Infof("Accrual Service : job %v added", job)
		return
	case <-time.After(AddTimeout):
		logger.Infof("Accrual Service : job %v can not be added", job)
		return
	}
}

func (as *AccrualService) PrepareNewJobQueue() {
	limit := QueueLength / 10 * 8
	if limit == 0 {
		limit = QueueLength - 2
	}
	orders, err := as.store.GetUnprocessedOrders(as.ctx, limit)
	// TODO Stopped here
	_ = orders
	_ = err
}

func (as *AccrualService) ProcessNewJobQueue() {
	const maxWorkers = 2

	var wg sync.WaitGroup
	logger.Info("Accrual Service: Queue processing started")
	d := time.Now().Add(8 * time.Second)                         // dev & debug only
	ctx, cancel := context.WithDeadline(context.Background(), d) // dev & debug only
	defer cancel()

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			logger.Infof("Accrual Service: Client goroutines #%d started", i)
			client, _ := NewAccrualClient(as.address)
			client.DoWork(ctx, as.newJobQueue)
		}(i)
	}
	wg.Wait()
	logger.Info("Accrual Service: Queue processing stopped")
}
