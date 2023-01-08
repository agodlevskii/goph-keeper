package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/models"
	"github.com/agodlevskii/goph-keeper/internal/pkg/services/data"
	"github.com/agodlevskii/goph-keeper/internal/pkg/services/text"
)

func TestNewTextService(t *testing.T) {
	ds := initDataMS(t)
	tests := []struct {
		name string
		want *TextService
	}{
		{
			name: "Service creation",
			want: &TextService{textMS: text.NewService(ds)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewTextService(ds))
		})
	}
}

func TestTextService_DeleteText(t *testing.T) {
	type args struct {
		uid string
		id  string
	}
	tests := []struct {
		name    string
		ds      data.Service
		repo    map[string]models.TextResponse
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
			repo:    map[string]models.TextResponse{"test1": {ID: "test1", UID: "test1", Name: "test1", Data: "test1"}},
			args:    args{uid: "test", id: "test"},
			wantErr: ErrTextNotFound,
		},
		{
			name: "Data is present and deleted",
			repo: map[string]models.TextResponse{"test": {ID: "test", UID: "test", Name: "test", Data: "test"}},
			args: args{uid: "test", id: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, ids := initTextService(t, tt.repo)
			if len(ids) > 0 {
				for id, v := range ids {
					if tt.args.id == id {
						tt.args.id = v.ID
					}
				}
			}

			err := s.DeleteText(context.Background(), tt.args.uid, tt.args.id)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestTextService_GetAllTexts(t *testing.T) {
	tests := []struct {
		name    string
		uid     string
		repo    map[string]models.TextResponse
		want    []models.TextResponse
		wantErr error
	}{
		{
			name:    "Missing UID",
			wantErr: ErrBadArguments,
		},
		{
			name: "No data",
			uid:  "test1",
			repo: map[string]models.TextResponse{"test": {ID: "test", UID: "test", Name: "test", Data: "test"}},
			want: []models.TextResponse{},
		},
		{
			name: "Data found",
			uid:  "test",
			repo: map[string]models.TextResponse{
				"test":  {ID: "test", UID: "test", Name: "test", Data: "test"},
				"test1": {ID: "test1", UID: "test1", Name: "test", Data: "test"},
			},
			want: []models.TextResponse{{UID: "test", Name: "test"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := initTextService(t, tt.repo)
			got, err := s.GetAllTexts(context.Background(), tt.uid)
			if len(got) == 0 {
				assert.Equal(t, tt.want, got)
			} else {
				assert.Equal(t, tt.want[0].Name, got[0].Name)
			}
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestTextService_GetTextByID(t *testing.T) {
	type args struct {
		uid string
		id  string
	}
	tests := []struct {
		name    string
		args    args
		repo    map[string]models.TextResponse
		want    models.TextResponse
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
			args:    args{uid: "test1", id: "test1"},
			repo:    map[string]models.TextResponse{"test": {ID: "test", UID: "test", Name: "test", Data: "test"}},
			wantErr: ErrTextNotFound,
		},
		{
			name: "Data found",
			args: args{uid: "test", id: "test"},
			repo: map[string]models.TextResponse{"test": {ID: "test", UID: "test", Name: "test", Data: "test"}},
			want: models.TextResponse{ID: "test", Name: "test", Data: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, ids := initTextService(t, tt.repo)
			if len(ids) > 0 {
				for id, v := range ids {
					if tt.args.id == id {
						tt.args.id = v.ID
						tt.want.ID = v.ID
					}
				}
			}

			got, err := s.GetTextByID(context.Background(), tt.args.uid, tt.args.id)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestHandler_StoreText(t *testing.T) {
	type args struct {
		uid  string
		text models.TextRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "Missing UID",
			wantErr: ErrBadArguments,
		},
		{
			name:    "Empty request",
			args:    args{uid: "test"},
			wantErr: ErrBadArguments,
		},
		{
			name:    "Empty name",
			args:    args{uid: "test", text: models.TextRequest{Data: "test"}},
			wantErr: ErrBadArguments,
		},
		{
			name:    "Empty data",
			args:    args{uid: "test", text: models.TextRequest{Name: "test"}},
			wantErr: ErrBadArguments,
		},
		{
			name: "Data saved",
			args: args{uid: "test", text: models.TextRequest{Name: "test", Data: "test"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := initTextService(t, nil)
			_, err := s.StoreText(context.Background(), tt.args.uid, tt.args.text)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func initTextService(t *testing.T,
	repo map[string]models.TextResponse,
) (*TextService, map[string]models.TextResponse) {
	s := TextService{textMS: text.NewService(initDataMS(t))}
	newRepo := make(map[string]models.TextResponse, len(repo))
	for iid, v := range repo {
		id, err := s.StoreText(context.Background(), v.UID, models.TextRequest{
			Name: v.Name,
			Data: v.Data,
			Note: v.Note,
		})
		if err != nil {
			t.Fatal(err)
		}
		v.ID = id
		newRepo[iid] = v
	}
	return &s, newRepo
}
