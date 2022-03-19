package models

import (
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"strings"
	"time"
)

type Money int64

func (a Money) MarshalJSON() ([]byte, error) {
	whole := strconv.FormatInt(int64(a)/100, 10)
	frac := ""
	rem := int64(a) % 100
	if rem > 0 {
		frac = "."
		if rem%10 > 0 {
			if rem/10 == 0 {
				frac += "0"
			}
			frac += strconv.FormatInt(rem, 10)
		} else {
			frac += strconv.FormatInt(rem/10, 10)
		}
	}
	return []byte(whole + frac), nil
}

func (a *Money) UnmarshalJSON(data []byte) error {
	s := string(data)
	switch strings.Count(s, ".") {
	case 0:
		res, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return errors.New("Money unmarshal error")
		}
		*a = Money(res * 100)
		return nil
	case 1:
		s += "00" // 123.4 -> 123.400; 123.45 -> 123.4500; 123. -> 123.00 (invalid json number, but still works)
		parts := strings.Split(s, ".")
		whole, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return errors.New("Money unmarshal error")
		}

		frac, err := strconv.ParseInt(parts[1][:2], 10, 64)
		if err != nil {
			return errors.New("Money unmarshal error")
		}
		*a = Money(whole*100 + frac)
		return nil
	default:
		return errors.New("Money unmarshal error")
	}
}

type AccrualResponse struct {
	OrderNumber string `json:"order"`
	Status      string `json:"status"`
	Accrual     Money  `json:"accrual"`
}

type Balance struct {
	Current   Money `json:"current"`
	Withdrawn Money `json:"withdrawn"`
}

// TODO check if ints will do instead of floats
// type BalanceResponse struct {
// 	Current   int64 `json:"current"`
// 	Withdrawn int64 `json:"withdrawn"`
// }

type WithdrawalInfo struct {
	OrderNumber string    `json:"order"`
	Sum         Money     `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

// TODO check if ints will do instead of floats
type WithdrawRequest struct {
	OrderNumber string `json:"order"`
	Sum         Money  `json:"sum"`
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
