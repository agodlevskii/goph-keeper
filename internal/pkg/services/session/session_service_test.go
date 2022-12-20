package session

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/agodlevskii/goph-keeper/internal/pkg/jwt"
)

func TestNewService(t *testing.T) {
	tests := []struct {
		name         string
		repoURL      string
		wantRepoType string
		wantErr      bool
	}{
		{
			name:         "Repo URL is missing",
			wantRepoType: "*session.BasicRepo",
		},
		{
			name:         "Wrong Repo URL is present",
			repoURL:      "postgres://localhost:5432/test",
			wantRepoType: "*session.DBRepo",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewService(tt.repoURL)
			assert.Equal(t, tt.wantErr, err != nil)

			rRepo := reflect.ValueOf(got.db)
			assert.Equal(t, tt.wantRepoType, rRepo.Type().String())
		})
	}
}

func TestService_DeleteSession(t *testing.T) {
	tests := []struct {
		name    string
		cid     string
		wantErr error
	}{
		{
			name:    "Missing client ID",
			wantErr: ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{db: initBasicRepo(nil)}
			err := s.DeleteSession(context.Background(), tt.cid)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestService_GenerateToken(t *testing.T) {
	tests := []struct {
		name    string
		uid     string
		wantLen int
		wantErr error
	}{
		{
			name:    "Empty ID",
			wantErr: ErrEmptyUID,
		},
		{
			name:    "Correct ID",
			uid:     "test-user",
			wantLen: 129,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{db: initBasicRepo(nil)}
			got, err := s.GenerateToken(tt.uid)
			assert.Equal(t, tt.wantLen, len(got))
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestService_GetUIDFromToken(t *testing.T) {
	token, err := jwt.EncodeToken("test-user", 0)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		token   string
		want    string
		wantErr error
	}{
		{
			name:    "Empty token",
			wantErr: ErrEmptyToken,
		},
		{
			name:  "Correct token",
			token: token,
			want:  "test-user",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{db: initBasicRepo(nil)}
			got, err := s.GetUIDFromToken(tt.token)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_RestoreSession(t *testing.T) {
	token, err := jwt.EncodeToken("test-user", 0)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		repo    map[string]string
		cid     string
		want    string
		wantErr error
	}{
		{
			name:    "Missing client ID",
			wantErr: ErrNotFound,
		},
		{
			name:    "Existing client ID, missing parameter",
			repo:    map[string]string{"testID": token},
			wantErr: ErrNotFound,
		},
		{
			name:    "Existing client ID, wrong parameter",
			repo:    map[string]string{"testID": token},
			cid:     "testID1",
			wantErr: ErrNotFound,
		},
		{
			name: "Existing client ID, correct parameter",
			repo: map[string]string{"testID": token},
			cid:  "testID",
			want: token,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{db: initBasicRepo(tt.repo)}
			got, err := s.RestoreSession(context.Background(), tt.cid)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestService_StoreSession(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantLen int
		wantErr error
	}{
		{
			name:    "Token is missing",
			wantErr: ErrIncorrectData,
			wantLen: 27,
		},
		{
			name:    "Token is present",
			token:   "Test token",
			wantLen: 27,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{db: initBasicRepo(nil)}

			got, err := s.StoreSession(context.Background(), tt.token)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantLen, len(got))
		})
	}
}

func Test_generateClientID(t *testing.T) {
	tests := []struct {
		name    string
		wantLen int
	}{
		{
			name:    "Token is generated",
			wantLen: 27,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateClientID()
			assert.Equal(t, tt.wantLen, len(got))
		})
	}
}

func TestService_IsTokenExpired(t *testing.T) {
	token, err := jwt.EncodeToken("test-expired1", 0)
	if err != nil {
		t.Fatal(err)
	}
	expToken, err := jwt.EncodeToken("test-expired2", time.Hour*-1)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		token   string
		want    bool
		wantErr error
	}{
		{
			name:    "Expired token",
			token:   expToken,
			want:    true,
			wantErr: ErrTokenExpired,
		},
		{
			name:  "Valid token",
			token: token,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{db: initBasicRepo(nil)}
			got, err := s.IsTokenExpired(tt.token)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
