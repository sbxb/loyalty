package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckLuhn(t *testing.T) {
	tests := []struct {
		num  string
		want bool
	}{
		{
			num:  "12345678903",
			want: true,
		},
		{
			num:  "49927398716",
			want: true,
		},
		{
			num:  "1234567812345670",
			want: true,
		},
		{
			num:  "49927398717",
			want: false,
		},
		{
			num:  "1234567812345678",
			want: false,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, CheckLuhn(tt.num), tt.want)
	}
}

func TestIsAllDigits(t *testing.T) {
	tests := []struct {
		num  string
		want bool
	}{
		{
			num:  "0",
			want: true,
		},
		{
			num:  "123",
			want: true,
		},
		{
			num:  "-99",
			want: false,
		},
		{
			num:  "1two3",
			want: false,
		},
		{
			num:  "0xcafe",
			want: false,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, isAllDigits(tt.num), tt.want)
	}
}
