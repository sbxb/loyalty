package models

import (
	"strings"
	"time"
)

type Order struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    Money     `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
}

const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessed  = "PROCESSED"
)

func (ord *Order) Validate() bool {
	ord.Number = strings.TrimSpace(ord.Number)
	if ord.Number == "" || !IsAllDigits(ord.Number) {
		return false
	}

	return true
}

// IsAllDigits tests if a string contains only digits
func IsAllDigits(str string) bool {
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
