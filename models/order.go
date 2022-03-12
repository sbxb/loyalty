package models

import (
	"strings"
	"time"
)

type Order struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    int64     `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
	// Exists      bool      `json:"-"`
	// Owner       string    `json:"-"`
}

const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessed  = "PROCESSED"
)

func (ord *Order) Validate() bool {
	ord.Number = strings.TrimSpace(ord.Number)
	if ord.Number == "" || !isAllDigits(ord.Number) {
		return false
	}

	return true
}

// isAllDigits tests if a string contains only digits
func isAllDigits(str string) bool {
	for _, r := range str {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// checkLuhn tests if a string containing only digits is a valid Luhn number
func CheckLuhn(str string) bool {
	sum := 0
	even := false

	for i := len(str) - 1; i >= 0; i-- {
		n := int(str[i] - '0')
		if even {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		even = !even
	}
	return sum%10 == 0
}
