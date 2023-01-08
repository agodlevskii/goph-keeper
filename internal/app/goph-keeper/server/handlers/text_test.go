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

func TestHandler_DeleteText(t *testing.T) {
	type args struct {
		uid string
		id  string
	}
	tests := []struct {
		name string
		repo map[string]models.TextResponse
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
			repo: map[string]models.TextResponse{"test": {ID: "test", UID: "test", Name: "test", Data: "test"}},
			args: args{uid: "test1", id: "test1"},
			want: httpRes{code: http.StatusNotFound},
		},
		{
			name: "Data is present and deleted",
			repo: map[string]models.TextResponse{"test": {ID: "test", UID: "test", Name: "test", Data: "test"}},
			args: args{uid: "test", id: "test"},
			want: httpRes{code: http.StatusOK},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts, ids := initTextService(t, tt.repo)
			if len(ids) > 0 {
				for id, v := range ids {
					if tt.args.id == id {
						tt.args.id = v.ID
					}
				}
			}

			h := Handler{textService: ts}
			r := initTestRequest(t, http.MethodDelete, textURL, tt.args.id, tt.args.uid, nil)
			w := httptest.NewRecorder()

			h.DeleteText()(w, r)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}

func TestHandler_GetAllTexts(t *testing.T) {
	tests := []struct {
		name string
		uid  string
		repo map[string]models.TextResponse
		want httpRes
	}{
		{
			name: "Missing UID",
			want: httpRes{code: http.StatusBadRequest},
		},
		{
			name: "No data",
			uid:  "test1",
			repo: map[string]models.TextResponse{"test": {ID: "test", UID: "test", Name: "test", Data: "test"}},
			want: httpRes{
				code: http.StatusOK,
				resp: "[]",
			},
		},
		{
			name: "Data found",
			uid:  "test",
			repo: map[string]models.TextResponse{
				"test":  {ID: "test", UID: "test", Name: "test", Data: "test"},
				"test1": {ID: "test1", UID: "test1", Name: "test", Data: "test"},
			},
			want: httpRes{
				code: http.StatusOK,
				resp: `[{UID: "test", Name: "test", Data: "test"}]`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts, _ := initTextService(t, tt.repo)
			h := Handler{textService: ts}
			r := initTestRequest(t, http.MethodGet, textURL, "", tt.uid, nil)
			w := httptest.NewRecorder()

			h.GetAllTexts()(w, r)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}

func TestHandler_GetTextByID(t *testing.T) {
	type args struct {
		uid string
		id  string
	}
	tests := []struct {
		name string
		args args
		repo map[string]models.TextResponse
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
			args: args{uid: "test1", id: "test1"},
			repo: map[string]models.TextResponse{"test": {ID: "test", UID: "test", Name: "test", Data: "test"}},
			want: httpRes{code: http.StatusNotFound},
		},
		{
			name: "Data found",
			args: args{uid: "test", id: "test"},
			repo: map[string]models.TextResponse{"test": {ID: "test", UID: "test", Name: "test", Data: "test"}},
			want: httpRes{
				code: http.StatusOK,
				resp: `{ID: "test", Name: "test", Data: "test"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts, ids := initTextService(t, tt.repo)
			if len(ids) > 0 {
				for id, v := range ids {
					if tt.args.id == id {
						tt.args.id = v.ID
						tt.want.resp = fmt.Sprintf("{ID: %s}", v.ID)
					}
				}
			}

			h := Handler{textService: ts}
			r := initTestRequest(t, http.MethodGet, cardURL, tt.args.id, tt.args.uid, nil)
			w := httptest.NewRecorder()

			h.GetTextByID()(w, r)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}

func TestHandler_StoreText(t *testing.T) {
	type args struct {
		uid string
	}
	type fields struct {
		req models.TextRequest
	}
	tests := []struct {
		name   string
		args   args
		fields fields
		want   httpRes
	}{
		{
			name: "Missing UID",
			want: httpRes{code: http.StatusBadRequest},
		},
		{
			name: "Empty request",
			args: args{uid: "test"},
			want: httpRes{code: http.StatusBadRequest},
		},
		{
			name:   "Empty name",
			args:   args{uid: "test"},
			fields: fields{req: models.TextRequest{Data: "test"}},
			want:   httpRes{code: http.StatusBadRequest},
		},
		{
			name:   "Empty data",
			args:   args{uid: "test"},
			fields: fields{req: models.TextRequest{Name: "test"}},
			want:   httpRes{code: http.StatusBadRequest},
		},
		{
			name:   "Data saved",
			args:   args{uid: "test"},
			fields: fields{req: models.TextRequest{Name: "test", Data: "test"}},
			want:   httpRes{code: http.StatusOK},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts, _ := initTextService(t, nil)
			h := Handler{textService: ts}
			r := initTestRequest(t, http.MethodPost, binaryURL, "", tt.args.uid, tt.fields.req)
			w := httptest.NewRecorder()

			h.StoreText()(w, r)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}

func initTextService(t *testing.T,
	repo map[string]models.TextResponse,
) (*services.TextService, map[string]models.TextResponse) {
	s := services.NewTextService(initDataMS(t))
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
	return s, newRepo
}
