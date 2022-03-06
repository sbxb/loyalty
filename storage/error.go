package storage

import "errors"

var ErrLoginAlreadyExists = errors.New("Login already exists")
var ErrLoginMissing = errors.New("Login missing")

var ErrOrderAlreadyExists = errors.New("Order already exists")

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
