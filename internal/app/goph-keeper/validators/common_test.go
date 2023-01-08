package validators

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItemName(t *testing.T) {
	tests := []struct {
		name    string
		item    string
		wantErr error
	}{
		{
			name:    "Short name",
			item:    "ts",
			wantErr: errors.New("the value must be at least 3 characters long"),
		},
		{
			name:    "Long name",
			item:    "gjsldfiodhgiouehautghoudshfguehsuilghui4hguiheiusgbhvuidhfughaujhgrfiuheaiufhjuegf",
			wantErr: errors.New("the value is limited to 50 characters"),
		},
		{
			name: "Correct name",
			item: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantErr, ItemName(tt.item))
		})
	}
}
