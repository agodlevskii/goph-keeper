package data

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
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
			wantRepoType: "*data.BasicRepo",
		},
		{
			name:         "Wrong Repo URL is present",
			repoURL:      "postgres://localhost:5432/test",
			wantRepoType: "*data.DBRepo",
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

func TestService_StoreSecureDataFromPayload(t *testing.T) {
	type args struct {
		uid     string
		payload any
		t       StorageType
	}
	tests := []struct {
		name    string
		repo    map[string]SecureData
		args    args
		want    string
		wantErr error
	}{
		{
			name:    "User ID is missing",
			args:    args{payload: []byte("test"), t: SCard},
			wantErr: ErrEmpty,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{db: initBasicRepo(tt.repo)}
			got, err := s.StoreSecureDataFromPayload(context.Background(), tt.args.uid, tt.args.payload, tt.args.t)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
