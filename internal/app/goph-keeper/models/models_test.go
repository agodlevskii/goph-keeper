package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinaryResponse_TableRow(t *testing.T) {
	tests := []struct {
		name  string
		model BinaryResponse
		want  []string
	}{
		{
			name: "Binary model",
			model: BinaryResponse{
				UID:  "testUID",
				ID:   "testID",
				Name: "testName",
				Note: "testNote",
			},
			want: []string{"testID", "testName", "", "testNote"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.model.TableRow())
		})
	}
}

func TestCardResponse_TableRow(t *testing.T) {
	tests := []struct {
		name  string
		model CardResponse
		want  []string
	}{
		{
			name: "Card model",
			model: CardResponse{
				UID:     "testUID",
				ID:      "testID",
				Name:    "testName",
				Number:  "123",
				Holder:  "test",
				ExpDate: "12/23",
				CVV:     "123",
				Note:    "testNote",
			},
			want: []string{"testID", "testName", "123", "test", "12/23", "123", "testNote"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.model.TableRow())
		})
	}
}

func TestPasswordResponse_TableRow(t *testing.T) {
	tests := []struct {
		name  string
		model PasswordResponse
		want  []string
	}{
		{
			name: "Password model",
			model: PasswordResponse{
				UID:      "testUID",
				ID:       "testID",
				Name:     "testName",
				User:     "testUser",
				Password: "testPassword",
				Note:     "testNote",
			},
			want: []string{"testID", "testName", "testUser", "testPassword", "testNote"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.model.TableRow())
		})
	}
}

func TestTextResponse_TableRow(t *testing.T) {
	tests := []struct {
		name  string
		model TextResponse
		want  []string
	}{
		{
			name: "Text model",
			model: TextResponse{
				UID:  "testUID",
				ID:   "testID",
				Name: "testName",
				Note: "testNote",
			},
			want: []string{"testID", "testName", "", "testNote"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.model.TableRow())
		})
	}
}
