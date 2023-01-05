package enc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name    string
		pass    string
		want    int
		wantErr error
	}{
		{
			name:    "Missing password",
			wantErr: ErrPasswordLength,
		},
		{
			name: "Correct password",
			pass: "test",
			want: 60,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, eErr := HashPassword(tt.pass)
			assert.Equal(t, tt.want, len(got))
			assert.Equal(t, tt.wantErr, eErr)
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	hp, err := HashPassword("test")
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		pwd  string
		hash string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Wrong password",
			args: args{pwd: "wrong", hash: hp},
		},
		{
			name: "Correct password",
			args: args{pwd: "test", hash: hp},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, VerifyPassword(tt.args.pwd, tt.args.hash))
		})
	}
}
