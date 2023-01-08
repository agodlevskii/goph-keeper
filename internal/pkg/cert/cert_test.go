package cert

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCertificatePool(t *testing.T) {
	tests := []struct {
		name     string
		fName    string
		path     string
		wantCert bool
		wantErr  bool
	}{
		{
			name:    "Missing path",
			wantErr: true,
		},
		{
			name:     "Correct path",
			fName:    "test.csv",
			path:     "test.csv",
			wantCert: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fName != "" {
				if _, err := os.Create(tt.fName); err != nil {
					t.Fatal(err)
				}
			}
			got, err := GetCertificatePool(tt.path)
			assert.Equal(t, tt.wantCert, got != nil)
			assert.Equal(t, tt.wantErr, err != nil)

			if tt.fName != "" {
				if err = os.Remove(tt.fName); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
