package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEncodeToken(t *testing.T) {
	type args struct {
		uid     string
		expTime time.Duration
	}
	type want struct {
		uid     string
		expired bool
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr error
	}{
		{
			name:    "Empty UID",
			wantErr: ErrTokenClaims,
		},
		{
			name: "Incorrect exp time",
			args: args{uid: "test", expTime: -1 * time.Minute},
			want: want{expired: true},
		},
		{
			name: "Presented UID",
			args: args{uid: "test"},
			want: want{uid: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeToken(tt.args.uid, tt.args.expTime)
			assert.Equal(t, tt.wantErr, err)

			if err == nil {
				valid, vErr := IsTokenExpired(got)
				assert.Equal(t, tt.want.expired, valid)

				if vErr == nil {
					uid, _ := GetUserIDFromToken(got)
					assert.Equal(t, tt.want.uid, uid)
				}
			}
		})
	}
}
