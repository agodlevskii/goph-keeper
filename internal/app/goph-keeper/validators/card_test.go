package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardCVV(t *testing.T) {
	tests := []struct {
		name    string
		cvv     string
		wantErr error
	}{
		{
			name:    "Missing CVV",
			wantErr: ErrCVVFormat,
		},
		{
			name:    "Short CVV",
			cvv:     "12",
			wantErr: ErrCVVFormat,
		},
		{
			name:    "Long CVV",
			cvv:     "1234",
			wantErr: ErrCVVFormat,
		},
		{
			name:    "Alphabetical CVV",
			cvv:     "abc",
			wantErr: ErrCVVFormat,
		},
		{
			name: "Correct CVV",
			cvv:  "123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantErr, CardCVV(tt.cvv))
		})
	}
}

func TestCardExpDate(t *testing.T) {
	tests := []struct {
		name    string
		date    string
		wantErr error
	}{
		{
			name:    "Missing date",
			wantErr: ErrExpDateFormat,
		},
		{
			name:    "Incorrectly formatted date",
			date:    "123/45",
			wantErr: ErrExpDateFormat,
		},
		{
			name:    "Date in the past",
			date:    "11/11",
			wantErr: ErrExpDateInPast,
		},
		{
			name: "Correct date",
			date: "12/68",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantErr, CardExpDate(tt.date))
		})
	}
}

func TestCardNumber(t *testing.T) {
	tests := []struct {
		name    string
		num     string
		wantErr error
	}{
		{
			name:    "Missing number",
			wantErr: ErrCardNumberFormat,
		},
		{
			name:    "Number is too long",
			num:     "1234 5678 9012 3456 7890",
			wantErr: ErrCardNumberFormat,
		},
		{
			name: "Correct number",
			num:  "1234567891023732",
		},
		{
			name: "Correct number with spaces",
			num:  "1234 5678 9012 3456",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantErr, CardNumber(tt.num))
		})
	}
}
