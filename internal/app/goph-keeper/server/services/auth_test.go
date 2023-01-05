package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/models"
	"github.com/agodlevskii/goph-keeper/internal/pkg/jwt"
	"github.com/agodlevskii/goph-keeper/internal/pkg/services/auth"
	"github.com/agodlevskii/goph-keeper/internal/pkg/services/session"
	"github.com/agodlevskii/goph-keeper/internal/pkg/services/user"
)

func TestNewAuthService(t *testing.T) {
	ss, us := initSessionUserMS(t)
	tests := []struct {
		name    string
		want    *AuthService
		wantErr bool
	}{
		{
			name: "Service creation",
			want: &AuthService{authMS: auth.NewService(ss, us)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as, err := NewAuthService("")
			assert.Equal(t, tt.want, as)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func initSessionUserMS(t *testing.T) (session.Service, user.Service) {
	ss, err := session.NewService("")
	if err != nil {
		t.Fatal(err)
	}
	us, err := user.NewService("")
	if err != nil {
		t.Fatal(err)
	}
	return ss, us
}

func TestAuthService_Register(t *testing.T) {
	type args struct {
		user models.UserRequest
		req  models.UserRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "Empty name",
			args:    args{req: models.UserRequest{Password: "test"}},
			wantErr: ErrBadArguments,
		},
		{
			name:    "Empty password",
			args:    args{req: models.UserRequest{Name: "test"}},
			wantErr: ErrBadArguments,
		},
		{
			name: "User exists",
			args: args{
				req:  models.UserRequest{Name: "test", Password: "test"},
				user: models.UserRequest{Name: "test", Password: "test1"},
			},
			wantErr: user.ErrExists,
		},
		{
			name: "User is registered",
			args: args{
				req:  models.UserRequest{Name: "test", Password: "test"},
				user: models.UserRequest{Name: "test1", Password: "test1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewAuthService("")
			if err != nil {
				t.Fatal(err)
			}
			if tt.args.user.Name != "" {
				if err = s.Register(context.Background(), tt.args.user); err != nil {
					t.Fatal(err)
				}
			}
			assert.Equal(t, tt.wantErr, s.Register(context.Background(), tt.args.req))
		})
	}
}

func TestAuthService_Authorize(t *testing.T) {
	token, err := jwt.EncodeToken("test_id", 0)
	if err != nil {
		t.Fatal(err)
	}
	expToken, err := jwt.EncodeToken("bad_id", -1*time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name:    "Missing token",
			wantErr: ErrBadArguments,
		},
		{
			name: "Valid token",
			args: args{token: token},
			want: "test_id",
		},
		{
			name:    "Expired token",
			args:    args{token: expToken},
			wantErr: auth.ErrSessionExpired,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, sErr := NewAuthService("")
			if sErr != nil {
				t.Fatal(sErr)
			}
			got, sErr := s.Authorize(tt.args.token)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, sErr)
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	type args struct {
		cid  string
		user models.UserRequest
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr error
	}{
		{
			name:    "Empty user name",
			args:    args{user: models.UserRequest{Password: "test"}},
			wantErr: ErrBadArguments,
		},
		{
			name:    "Empty user password",
			args:    args{user: models.UserRequest{Name: "test"}},
			wantErr: ErrBadArguments,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, sErr := NewAuthService("")
			if sErr != nil {
				t.Fatal(sErr)
			}
			got, got1, err := s.Login(context.Background(), tt.args.cid, tt.args.user)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
