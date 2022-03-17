package models

import (
	"fmt"
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
		assert.Equal(t, IsAllDigits(tt.num), tt.want)
	}
}

func TestMakeLuhnNumbers(t *testing.T) {
	t.Skip()
	min := 1100
	max := 9999
	for i := min; i <= max; i++ {
		res := CheckLuhn(fmt.Sprintf("%d", i))
		if res {
			t.Log(i)
		}
	}
}

// 1149 1156 1172
// 2238 2253 2279
// 3327 3376 3384
// 4416 4457 4481
// 5512 5538 5587
// 6619 6635 6684
// 7724 7740 7781
// 8821 8862 8888
// 9936 9944 9977
