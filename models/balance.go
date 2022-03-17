package models

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
	"time"
)

type AccrualResponse struct {
	OrderNumber string `json:"order"`
	Status      string `json:"status"`
	Accrual     int64  `json:"accrual"`
}

type Balance struct {
	Current   int64
	Withdrawn int64
}

// TODO check if ints will do instead of floats
type BalanceResponse struct {
	Current   int64 `json:"current"`
	Withdrawn int64 `json:"withdrawn"`
}

type WithdrawalInfo struct {
	OrderNumber string    `json:"order"`
	Sum         int64     `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

// TODO check if ints will do instead of floats
type WithdrawRequest struct {
	OrderNumber string `json:"order"`
	Sum         int64  `json:"sum"`
}

func ReadWithdrawRequestFromBody(r io.Reader) (*WithdrawRequest, error) {
	req := &WithdrawRequest{}

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	if err := dec.Decode(req); err != nil {
		return nil, errors.New("Bad request: " + err.Error())
	}

	return req, nil
}

func (req *WithdrawRequest) Validate() bool {
	req.OrderNumber = strings.TrimSpace(req.OrderNumber)
	if req.OrderNumber == "" || !IsAllDigits(req.OrderNumber) {
		return false
	}
	if !CheckLuhn(req.OrderNumber) {
		return false
	}
	return true
}
