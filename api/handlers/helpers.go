package handlers

import (
	"io"
	"net/http"

	"github.com/sbxb/loyalty/internal/logger"
	"github.com/sbxb/loyalty/models"
)

type OrderPostError struct {
	msg  string
	Code int
}

func (ope OrderPostError) Error() string {
	return ope.msg
}

func NewOrderPostError(msg string, code int) *OrderPostError {
	return &OrderPostError{
		msg:  msg,
		Code: code,
	}
}

func ReadOrderNumberFromBody(body io.ReadCloser) (*models.Order, *OrderPostError) {
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, NewOrderPostError("Server failed to read the request's body", http.StatusInternalServerError)
	}

	order := &models.Order{Number: string(data)}
	logger.Infof("%v\n", *order)
	if !order.Validate() {
		return nil, NewOrderPostError("wrong request format", http.StatusBadRequest)
	}

	if !models.CheckLuhn(order.Number) {
		return nil, NewOrderPostError("wrong number format", http.StatusUnprocessableEntity)
	}

	return order, nil
}
