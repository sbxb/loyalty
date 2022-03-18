package accrual

import (
	"context"

	"github.com/sbxb/loyalty/internal/logger"
	"github.com/sbxb/loyalty/models"
	"github.com/sbxb/loyalty/storage"
)

type SimpleAccrualService struct {
	store   storage.Storage
	client  *AccrualClient
	address string
}

func NewSimpleAccrualService(st storage.Storage, address string) *SimpleAccrualService {
	logger.Info("Accrual Service : created")
	client, _ := NewAccrualClient(address)
	return &SimpleAccrualService{
		store:   st,
		client:  client,
		address: address,
	}
}

func (sas *SimpleAccrualService) DoAccrualStuff(orderNumber string) bool {
	var retry bool
	logger.Infof("DoAccrualStuff(): Trying to get " + sas.client.url + orderNumber)
	resp, err := sas.client.client.Get(sas.client.url + orderNumber)
	if err != nil {
		logger.Infof("DoAccrualStuff(): Get request failed")
		return retry
	}
	ar, clientErr := processResponse(resp)
	if clientErr != nil {
		logger.Warning("DoAccrualStuff: " + clientErr.Error())
		return retry
	}
	logger.Infof("DoAccrualStuff(): got response %v", ar)
	switch ar.Status {
	case "REGISTERED":
		logger.Warning("DoAccrualStuff: Accrual Server has not processed order yet, need another try")
		retry = true
	case models.OrderStatusInvalid:
		err := sas.store.UpdateOrderStatus(context.Background(), ar)
		if err != nil {
			logger.Warning("DoAccrualStuff: Store failed to change the order status")
		}
	case models.OrderStatusProcessing:
		err := sas.store.UpdateOrderStatus(context.Background(), ar)
		if err != nil {
			logger.Warning("DoAccrualStuff: Store failed to change the order status")
		}
		logger.Warning("DoAccrualStuff: Accrual Server has not processed order yet, need another try")
		retry = true
	case models.OrderStatusProcessed:
		err := sas.store.ProcessOrder(context.Background(), ar)
		if err != nil {
			logger.Warning("DoAccrualStuff: Store failed to process the order")
		}
	default:
		logger.Warning("DoAccrualStuff: Should never see this")
	}
	return retry
}
