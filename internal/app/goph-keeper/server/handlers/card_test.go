package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/models"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/services"
)

func TestHandler_DeleteCard(t *testing.T) {
	type args struct {
		uid string
		id  string
	}
	tests := []struct {
		name string
		repo map[string]models.CardResponse
		args args
		want httpRes
	}{
		{
			name: "UID is missing",
			args: args{id: "test"},
			want: httpRes{code: http.StatusBadRequest},
		},
		{
			name: "ID is missing",
			args: args{uid: "testID"},
			want: httpRes{code: http.StatusBadRequest},
		},
		{
			name: "Data is not present",
			repo: map[string]models.CardResponse{"test1": {ID: "test1", UID: "test1"}},
			args: args{uid: "test", id: "test"},
			want: httpRes{code: http.StatusNotFound},
		},
		{
			name: "Data is present and deleted",
			repo: map[string]models.CardResponse{"test": {ID: "test", UID: "test"}},
			args: args{uid: "test", id: "test"},
			want: httpRes{code: http.StatusOK},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs, ids := initCardService(t, tt.repo)
			if len(ids) > 0 {
				for id, v := range ids {
					if tt.args.id == id {
						tt.args.id = v.ID
					}
				}
			}

			h := Handler{cardService: cs}
			r := initTestRequest(t, http.MethodDelete, cardURL, tt.args.id, tt.args.uid, nil)
			w := httptest.NewRecorder()

			h.DeleteCard()(w, r)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}

func TestHandler_GetAllCards(t *testing.T) {
	tests := []struct {
		name string
		uid  string
		repo map[string]models.CardResponse
		want httpRes
	}{
		{
			name: "Missing UID",
			want: httpRes{code: http.StatusBadRequest},
		},
		{
			name: "No data",
			uid:  "test1",
			repo: map[string]models.CardResponse{"test": {UID: "test", ID: "test"}},
			want: httpRes{
				code: http.StatusOK,
				resp: "[]",
			},
		},
		{
			name: "Data found",
			uid:  "test",
			repo: map[string]models.CardResponse{
				"test":  {UID: "test", Name: "test"},
				"test1": {UID: "test1", Name: "test1"},
			},
			want: httpRes{
				code: http.StatusOK,
				resp: `[{UID: "test", Name: "test"}]`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs, _ := initCardService(t, tt.repo)
			h := Handler{cardService: cs}
			r := initTestRequest(t, http.MethodGet, cardURL, "", tt.uid, nil)
			w := httptest.NewRecorder()

			h.GetAllCards()(w, r)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}

func TestHandler_GetCardByID(t *testing.T) {
	type args struct {
		uid string
		id  string
	}
	tests := []struct {
		name string
		args args
		repo map[string]models.CardResponse
		want httpRes
	}{
		{
			name: "Missing arguments",
			want: httpRes{code: http.StatusBadRequest},
		},
		{
			name: "Missing UID",
			args: args{id: "test"},
			want: httpRes{code: http.StatusBadRequest},
		},
		{
			name: "Missing ID",
			args: args{uid: "test"},
			want: httpRes{code: http.StatusBadRequest},
		},
		{
			name: "No data",
			args: args{uid: "test", id: "test"},
			repo: map[string]models.CardResponse{"test1": {UID: "test1", ID: "test1"}},
			want: httpRes{code: http.StatusNotFound},
		},
		{
			name: "Data found",
			args: args{uid: "test", id: "test"},
			repo: map[string]models.CardResponse{"test": {ID: "test", UID: "test"}},
			want: httpRes{
				code: http.StatusOK,
				resp: `{ID: "test"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs, ids := initCardService(t, tt.repo)
			if len(ids) > 0 {
				for id, v := range ids {
					if tt.args.id == id {
						tt.args.id = v.ID
						tt.want.resp = fmt.Sprintf("{ID: %s}", v.ID)
					}
				}
			}

			h := Handler{cardService: cs}
			r := initTestRequest(t, http.MethodGet, cardURL, tt.args.id, tt.args.uid, nil)
			w := httptest.NewRecorder()

			h.GetCardByID()(w, r)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}

func initCardService(t *testing.T,
	repo map[string]models.CardResponse,
) (*services.CardService, map[string]models.CardResponse) {
	s := services.NewCardService(initDataMS(t))
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
	return s, newRepo
}
