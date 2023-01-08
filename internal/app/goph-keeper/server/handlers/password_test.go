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

func TestHandler_DeletePassword(t *testing.T) {
	type args struct {
		uid string
		id  string
	}
	tests := []struct {
		name string
		repo map[string]models.PasswordResponse
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
			repo: map[string]models.PasswordResponse{"test1": {ID: "test1", UID: "test1"}},
			args: args{uid: "test", id: "test"},
			want: httpRes{code: http.StatusNotFound},
		},
		{
			name: "Data is present and deleted",
			repo: map[string]models.PasswordResponse{"test": {ID: "test", UID: "test"}},
			args: args{uid: "test", id: "test"},
			want: httpRes{code: http.StatusOK},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps, ids := initPasswordService(t, tt.repo)
			if len(ids) > 0 {
				for id, v := range ids {
					if tt.args.id == id {
						tt.args.id = v.ID
					}
				}
			}

			h := Handler{passwordService: ps}
			r := initTestRequest(t, http.MethodDelete, pStorageURL, tt.args.id, tt.args.uid, nil)
			w := httptest.NewRecorder()

			h.DeletePassword()(w, r)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}

func TestHandler_GetAllPasswords(t *testing.T) {
	tests := []struct {
		name string
		uid  string
		repo map[string]models.PasswordResponse
		want httpRes
	}{
		{
			name: "Missing UID",
			want: httpRes{code: http.StatusBadRequest},
		},
		{
			name: "No data",
			uid:  "test1",
			repo: map[string]models.PasswordResponse{"test": {UID: "test", ID: "test"}},
			want: httpRes{
				code: http.StatusOK,
				resp: "[]",
			},
		},
		{
			name: "Data found",
			uid:  "test",
			repo: map[string]models.PasswordResponse{
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
			ps, _ := initPasswordService(t, tt.repo)
			h := Handler{passwordService: ps}
			r := initTestRequest(t, http.MethodGet, pStorageURL, "", tt.uid, nil)
			w := httptest.NewRecorder()

			h.GetAllPasswords()(w, r)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}

func TestHandler_GetPasswordByID(t *testing.T) {
	type args struct {
		uid string
		id  string
	}
	tests := []struct {
		name string
		args args
		repo map[string]models.PasswordResponse
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
			repo: map[string]models.PasswordResponse{"test1": {UID: "test1", ID: "test1"}},
			want: httpRes{code: http.StatusNotFound},
		},
		{
			name: "Data found",
			args: args{uid: "test", id: "test"},
			repo: map[string]models.PasswordResponse{"test": {ID: "test", UID: "test"}},
			want: httpRes{
				code: http.StatusOK,
				resp: `{ID: "test"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps, ids := initPasswordService(t, tt.repo)
			if len(ids) > 0 {
				for id, v := range ids {
					if tt.args.id == id {
						tt.args.id = v.ID
						tt.want.resp = fmt.Sprintf("{ID: %s}", v.ID)
					}
				}
			}

			h := Handler{passwordService: ps}
			r := initTestRequest(t, http.MethodGet, cardURL, tt.args.id, tt.args.uid, nil)
			w := httptest.NewRecorder()

			h.GetPasswordByID()(w, r)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}

func initPasswordService(t *testing.T,
	repo map[string]models.PasswordResponse,
) (*services.PasswordService, map[string]models.PasswordResponse) {
	s := services.NewPasswordService(initDataMS(t))
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
	return s, newRepo
}
