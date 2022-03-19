package storage

import "errors"

var ErrLoginAlreadyExists = errors.New("login already exists")
var ErrLoginMissing = errors.New("login missing")

var ErrInsufficientFunds = errors.New("insufficient amount of loyalty points to withdraw")

var ErrOrderAlreadyExists = errors.New("order already exists")

type ExistingOrderError struct {
	Err    error
	UserID int
}

func NewExistingOrderError(userID int) *ExistingOrderError {
	return &ExistingOrderError{
		Err:    ErrOrderAlreadyExists,
		UserID: userID,
	}
}

func (eoe *ExistingOrderError) Error() string {
	return eoe.Err.Error()
}
