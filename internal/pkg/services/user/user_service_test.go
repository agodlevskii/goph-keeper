package user

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/agodlevskii/goph-keeper/internal/pkg/enc"
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
			wantRepoType: "*user.BasicRepo",
		},
		{
			name:         "Wrong Repo URL is present",
			repoURL:      "postgres://localhost:5432/test",
			wantRepoType: "*user.DBRepo",
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

func TestService_AddUser(t *testing.T) {
	tp, err := enc.HashPassword("test")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		user    User
		repo    map[string]User
		wantErr error
	}{
		{
			name:    "User is empty",
			wantErr: ErrCredMissing,
		},
		{
			name:    "User name is missing",
			user:    User{Password: "test"},
			wantErr: ErrCredMissing,
		},
		{
			name:    "User password is missing",
			user:    User{Name: "test"},
			wantErr: ErrCredMissing,
		},
		{
			name:    "User already exists",
			repo:    map[string]User{"test": {ID: "test", Name: "test", Password: tp}},
			user:    User{ID: "test", Name: "test", Password: "test"},
			wantErr: ErrExists,
		},
		{
			name: "User is present",
			repo: map[string]User{"test1": {ID: "test1", Name: "test1", Password: tp}},
			user: User{Name: "test", Password: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{db: initBasicRepo(tt.repo)}
			err = s.AddUser(context.Background(), tt.user)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
