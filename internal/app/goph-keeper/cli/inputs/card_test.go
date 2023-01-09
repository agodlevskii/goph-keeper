package inputs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardCVV(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "No input",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CardCVV()
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestCardExpDate(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "No input",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CardExpDate()
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestCardHolder(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "No input",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CardHolder()
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
