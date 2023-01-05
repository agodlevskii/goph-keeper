package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/models"
	"github.com/agodlevskii/goph-keeper/internal/pkg/services/data"
	"github.com/agodlevskii/goph-keeper/internal/pkg/services/password"
)

func TestNewPasswordService(t *testing.T) {
	ds := initDataMS(t)
	tests := []struct {
		name string
		want *PasswordService
	}{
		{
			name: "Service creation",
			want: &PasswordService{passwordMS: password.NewService(ds)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewPasswordService(ds))
		})
	}
}

func TestPasswordService_DeletePassword(t *testing.T) {
	type args struct {
		uid string
		id  string
	}
	tests := []struct {
		name    string
		ds      data.Service
		repo    map[string]models.PasswordResponse
		args    args
		wantErr error
	}{
		{
			name:    "Arguments are empty",
			wantErr: ErrBadArguments,
		},
		{
			name:    "ID is empty",
			args:    args{uid: "test"},
			wantErr: ErrBadArguments,
		},
		{
			name:    "User ID is empty",
			args:    args{id: "test"},
			wantErr: ErrBadArguments,
		},
		{
			name:    "Data is not present",
			repo:    map[string]models.PasswordResponse{"test1": {ID: "test1", UID: "test1"}},
			args:    args{uid: "test", id: "test"},
			wantErr: ErrPasswordNotFound,
		},
		{
			name: "Data is present and deleted",
			repo: map[string]models.PasswordResponse{"test": {ID: "test", UID: "test"}},
			args: args{uid: "test", id: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, ids := initPasswordService(t, tt.repo)
			if len(ids) > 0 {
				for id, v := range ids {
					if tt.args.id == id {
						tt.args.id = v.ID
					}
				}
			}

			err := s.DeletePassword(context.Background(), tt.args.uid, tt.args.id)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestPasswordService_GetAllPasswords(t *testing.T) {
	tests := []struct {
		name    string
		uid     string
		repo    map[string]models.PasswordResponse
		want    []models.PasswordResponse
		wantErr error
	}{
		{
			name:    "Missing UID",
			wantErr: ErrBadArguments,
		},
		{
			name: "No data",
			uid:  "test1",
			repo: map[string]models.PasswordResponse{"test": {UID: "test", ID: "test"}},
			want: []models.PasswordResponse{},
		},
		{
			name: "Data found",
			uid:  "test",
			repo: map[string]models.PasswordResponse{
				"test":  {UID: "test", Name: "test"},
				"test1": {UID: "test1", Name: "test1"},
			},
			want: []models.PasswordResponse{{UID: "test", Name: "test"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := initPasswordService(t, tt.repo)
			got, err := s.GetAllPasswords(context.Background(), tt.uid)
			if len(got) == 0 {
				assert.Equal(t, tt.want, got)
			} else {
				assert.Equal(t, tt.want[0].Name, got[0].Name)
			}
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestPasswordService_GetPasswordByID(t *testing.T) {
	type args struct {
		uid string
		id  string
	}
	tests := []struct {
		name    string
		args    args
		repo    map[string]models.PasswordResponse
		want    models.PasswordResponse
		wantErr error
	}{
		{
			name:    "Missing arguments",
			wantErr: ErrBadArguments,
		},
		{
			name:    "Missing UID",
			args:    args{id: "test"},
			wantErr: ErrBadArguments,
		},
		{
			name:    "Missing ID",
			args:    args{uid: "test"},
			wantErr: ErrBadArguments,
		},
		{
			name:    "No data",
			args:    args{uid: "test", id: "test"},
			repo:    map[string]models.PasswordResponse{"test1": {UID: "test1", ID: "test1"}},
			wantErr: ErrPasswordNotFound,
		},
		{
			name: "Data found",
			args: args{uid: "test", id: "test"},
			repo: map[string]models.PasswordResponse{"test": {ID: "test", UID: "test"}},
			want: models.PasswordResponse{ID: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, ids := initPasswordService(t, tt.repo)
			if len(ids) > 0 {
				for id, v := range ids {
					if tt.args.id == id {
						tt.args.id = v.ID
						tt.want.ID = v.ID
					}
				}
			}

			got, err := s.GetPasswordByID(context.Background(), tt.args.uid, tt.args.id)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func initPasswordService(t *testing.T,
	repo map[string]models.PasswordResponse,
) (*PasswordService, map[string]models.PasswordResponse) {
	s := PasswordService{passwordMS: password.NewService(initDataMS(t))}
	newRepo := make(map[string]models.PasswordResponse, len(repo))
	for iid, v := range repo {
		id, err := s.StorePassword(context.Background(), v.UID, models.PasswordRequest{
			Name:     v.Name,
			Password: v.Password,
			Note:     v.Note,
		})
		if err != nil {
			t.Fatal(err)
		}
		v.ID = id
		newRepo[iid] = v
	}
	return &s, newRepo
}
