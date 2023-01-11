package cert

import (
	"crypto/rand"
	"math/big"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCertificatePool(t *testing.T) {
	tests := []struct {
		name     string
		fName    string
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
			wantCert: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fName := generateCertFileName(t, tt.fName)
			if fName != "" {
				f, err := os.Create(fName)
				if err != nil {
					t.Fatal(err)
				}
				defer f.Close()
			}
			got, err := GetCertificatePool(fName)
			assert.Equal(t, tt.wantCert, got != nil)
			assert.Equal(t, tt.wantErr, err != nil)

			if fName != "" {
				if err = os.Remove(fName); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func generateCertFileName(t *testing.T, fName string) string {
	if fName == "" {
		return ""
	}

	r, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		t.Fatal(err)
	}
	return strconv.Itoa(int(r.Int64())) + fName
}
