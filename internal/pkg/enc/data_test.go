package enc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecryptData(t *testing.T) {
	enc, err := EncryptData([]byte("right_data_passed"))
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name    string
		data    []byte
		want    []byte
		wantErr error
	}{
		{
			name:    "Empty data",
			wantErr: ErrDataLength,
		},
		{
			name:    "Wrong encrypted data",
			data:    []byte("wrong_data_passed"),
			wantErr: ErrDecryption,
		},
		{
			name: "Right encrypted data",
			data: enc,
			want: []byte("right_data_passed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, dErr := DecryptData(tt.data)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, dErr)
		})
	}
}

func TestEncryptData(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    int
		wantErr error
	}{
		{
			name:    "Empty data",
			wantErr: ErrDataLength,
		},
		{
			name: "Correct data",
			data: []byte("right_data_passed"),
			want: 45,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, eErr := EncryptData(tt.data)
			assert.Equal(t, tt.want, len(got))
			assert.Equal(t, tt.wantErr, eErr)
		})
	}
}
