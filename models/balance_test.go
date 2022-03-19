package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TSNumber struct {
	Accrual Money `json:"accrual"`
}

func TestConvertToJSON(t *testing.T) {
	tests := []struct {
		json  string
		cents Money
	}{
		{
			json:  `{"accrual":1.99}`,
			cents: 199,
		},
		{
			json:  `{"accrual":2}`,
			cents: 200,
		},
		{
			json:  `{"accrual":2.01}`,
			cents: 201,
		},
		{
			json:  `{"accrual":2.1}`,
			cents: 210,
		},
		{
			json:  `{"accrual":2.11}`,
			cents: 211,
		},
		{
			json:  `{"accrual":2.5}`,
			cents: 250,
		},
		{
			json:  `{"accrual":123456789.99}`,
			cents: 12345678999,
		},
		{
			json:  `{"accrual":123456790}`,
			cents: 12345679000,
		},
		{
			json:  `{"accrual":123456790.01}`,
			cents: 12345679001,
		},
	}
	for _, tt := range tests {
		tj := TSNumber{tt.cents}
		jc, err := json.Marshal(tj)
		if err != nil {
			t.Log(err)
		}
		//t.Log(string(jc))
		assert.Equal(t, string(jc), tt.json)
	}
}

func TestConvertFromJSON(t *testing.T) {
	tests := []struct {
		json  string
		cents Money
	}{
		{
			json:  `{"accrual": 1.99}`,
			cents: 199,
		},
		{
			json:  `{"accrual": 2}`,
			cents: 200,
		},
		{
			json:  `{"accrual": 2.0}`,
			cents: 200,
		},
		{
			json:  `{"accrual": 2.00}`,
			cents: 200,
		},
		{
			json:  `{"accrual": 2.01}`,
			cents: 201,
		},
		{
			json:  `{"accrual": 2.5}`,
			cents: 250,
		},
		{
			json:  `{"accrual": 2.50}`,
			cents: 250,
		},
		{
			json:  `{"accrual": 123456789.99}`,
			cents: 12345678999,
		},
		{
			json:  `{"accrual": 123456790.00}`,
			cents: 12345679000,
		},
		{
			json:  `{"accrual": 123456790.01}`,
			cents: 12345679001,
		},
	}

	for _, tt := range tests {
		var res TSNumber
		err := json.Unmarshal([]byte(tt.json), &res)
		if err != nil {
			t.Log(err)
		}
		assert.Equal(t, res.Accrual, tt.cents)
	}
}
