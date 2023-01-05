package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/models"
	"github.com/agodlevskii/goph-keeper/internal/pkg/services/card"
)

func TestCardService_DeleteCard(t *testing.T) {
	type args struct {
		uid string
		id  string
	}
	tests := []struct {
		name    string
		repo    map[string]models.CardResponse
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
			repo:    map[string]models.CardResponse{"test1": {ID: "test1", UID: "test1"}},
			args:    args{uid: "test", id: "test"},
			wantErr: ErrCardNotFound,
		},
		{
			name: "Data is present and deleted",
			repo: map[string]models.CardResponse{"test": {ID: "test", UID: "test"}},
			args: args{uid: "test", id: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, ids := initCardService(t, tt.repo)
			if len(ids) > 0 {
				for id, v := range ids {
					if tt.args.id == id {
						tt.args.id = v.ID
					}
				}
			}

			err := s.DeleteCard(context.Background(), tt.args.uid, tt.args.id)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestCardService_GetAllCards(t *testing.T) {
	tests := []struct {
		name    string
		uid     string
		repo    map[string]models.CardResponse
		want    []models.CardResponse
		wantErr error
	}{
		{
			name:    "Missing UID",
			wantErr: ErrBadArguments,
		},
		{
			name: "No data",
			uid:  "test1",
			repo: map[string]models.CardResponse{"test": {UID: "test", ID: "test"}},
			want: []models.CardResponse{},
		},
		{
			name: "Data found",
			uid:  "test",
			repo: map[string]models.CardResponse{
				"test":  {UID: "test", Name: "test"},
				"test1": {UID: "test1", Name: "test1"},
			},
			want: []models.CardResponse{{UID: "test", Name: "test"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := initCardService(t, tt.repo)
			got, err := s.GetAllCards(context.Background(), tt.uid)
			if len(got) == 0 {
				assert.Equal(t, tt.want, got)
			} else {
				assert.Equal(t, tt.want[0].Name, got[0].Name)
			}
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestCardService_GetCardByID(t *testing.T) {
	type args struct {
		uid string
		id  string
	}
	tests := []struct {
		name    string
		args    args
		repo    map[string]models.CardResponse
		want    models.CardResponse
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
			repo:    map[string]models.CardResponse{"test1": {UID: "test1", ID: "test1"}},
			wantErr: ErrCardNotFound,
		},
		{
			name: "Data found",
			args: args{uid: "test", id: "test"},
			repo: map[string]models.CardResponse{"test": {ID: "test", UID: "test"}},
			want: models.CardResponse{ID: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, ids := initCardService(t, tt.repo)
			if len(ids) > 0 {
				for id, v := range ids {
					if tt.args.id == id {
						tt.args.id = v.ID
						tt.want.ID = v.ID
					}
				}
			}

			got, err := s.GetCardByID(context.Background(), tt.args.uid, tt.args.id)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestNewCardService(t *testing.T) {
	ds := initDataMS(t)
	tests := []struct {
		name string
		want *CardService
	}{
		{
			name: "Service creation",
			want: &CardService{cardMS: card.NewService(ds)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewCardService(ds))
		})
	}
}

func initCardService(t *testing.T,
	repo map[string]models.CardResponse,
) (*CardService, map[string]models.CardResponse) {
	s := CardService{cardMS: card.NewService(initDataMS(t))}
	newRepo := make(map[string]models.CardResponse, len(repo))
	for iid, v := range repo {
		id, err := s.StoreCard(context.Background(), v.UID, models.CardRequest{
			Name: v.Name,
			CVV:  v.CVV,
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
