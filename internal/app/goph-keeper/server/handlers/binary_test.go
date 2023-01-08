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
	"github.com/agodlevskii/goph-keeper/internal/pkg/services/data"
)

func TestHandler_DeleteBinary(t *testing.T) {
	type args struct {
		uid string
		id  string
	}
	tests := []struct {
		name string
		args args
		repo map[string]models.BinaryResponse
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
			repo: map[string]models.BinaryResponse{"t1": {UID: "t1", ID: "t1", Name: "t1", Data: []byte("t1")}},
			args: args{uid: "test", id: "test"},
			want: httpRes{code: http.StatusNotFound},
		},
		{
			name: "Data is present and deleted",
			repo: map[string]models.BinaryResponse{"t1": {UID: "t1", ID: "t1", Name: "t1", Data: []byte("t1")}},
			args: args{uid: "t1", id: "t1"},
			want: httpRes{code: http.StatusOK},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs, ids := initBinaryService(t, tt.repo)
			if len(ids) > 0 {
				for id, v := range ids {
					if tt.args.id == id {
						tt.args.id = v.ID
					}
				}
			}

			h := Handler{binaryService: bs}
			r := initTestRequest(t, http.MethodDelete, binaryURL, tt.args.id, tt.args.uid, nil)
			w := httptest.NewRecorder()

			h.DeleteBinary()(w, r)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}

func TestHandler_GetAllBinaries(t *testing.T) {
	tests := []struct {
		name string
		uid  string
		repo map[string]models.BinaryResponse
		want httpRes
	}{
		{
			name: "Missing UID",
			want: httpRes{code: http.StatusBadRequest},
		},
		{
			name: "No data",
			uid:  "test1",
			repo: map[string]models.BinaryResponse{"t1": {UID: "t1", ID: "t1", Name: "t1", Data: []byte("t1")}},
			want: httpRes{
				code: http.StatusOK,
				resp: "[]",
			},
		},
		{
			name: "Data found",
			uid:  "test",
			repo: map[string]models.BinaryResponse{
				"t":  {UID: "t", ID: "t", Name: "t", Data: []byte("t")},
				"t1": {UID: "t1", ID: "t1", Name: "t1", Data: []byte("t1")},
			},
			want: httpRes{
				code: http.StatusOK,
				resp: `[{UID: "test", Name: "test"}]`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs, _ := initBinaryService(t, tt.repo)
			h := Handler{binaryService: bs}
			r := initTestRequest(t, http.MethodGet, binaryURL, "", tt.uid, nil)
			w := httptest.NewRecorder()

			h.GetAllBinaries()(w, r)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}

func TestHandler_GetBinaryByID(t *testing.T) {
	type args struct {
		uid string
		id  string
	}
	tests := []struct {
		name string
		args args
		repo map[string]models.BinaryResponse
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
			repo: map[string]models.BinaryResponse{"test1": {UID: "t", ID: "t", Name: "t", Data: []byte("t")}},
			want: httpRes{code: http.StatusNotFound},
		},
		{
			name: "Data found",
			args: args{uid: "t", id: "t"},
			repo: map[string]models.BinaryResponse{"t": {UID: "t", ID: "t", Name: "t", Data: []byte("t")}},
			want: httpRes{
				code: http.StatusOK,
				resp: `{ID: "t", Name: "t"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs, ids := initBinaryService(t, tt.repo)
			if len(ids) > 0 {
				for id, v := range ids {
					if tt.args.id == id {
						tt.args.id = v.ID
						tt.want.resp = fmt.Sprintf("{ID: %s}", v.ID)
					}
				}
			}

			h := Handler{binaryService: bs}
			r := initTestRequest(t, http.MethodGet, binaryURL, tt.args.id, tt.args.uid, nil)
			w := httptest.NewRecorder()

			h.GetBinaryByID()(w, r)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}

func TestHandler_StoreBinary(t *testing.T) {
	type args struct {
		uid string
	}
	type fields struct {
		req models.BinaryRequest
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
			fields: fields{req: models.BinaryRequest{Data: []byte("test")}},
			want:   httpRes{code: http.StatusBadRequest},
		},
		{
			name:   "Empty data",
			args:   args{uid: "test"},
			fields: fields{req: models.BinaryRequest{Name: "test"}},
			want:   httpRes{code: http.StatusBadRequest},
		},
		{
			name:   "Data saved",
			args:   args{uid: "test"},
			fields: fields{req: models.BinaryRequest{Name: "test", Data: []byte("test")}},
			want:   httpRes{code: http.StatusOK},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs, _ := initBinaryService(t, nil)
			h := Handler{binaryService: bs}
			r := initTestRequest(t, http.MethodPost, binaryURL, "", tt.args.uid, tt.fields.req)
			w := httptest.NewRecorder()

			h.StoreBinary()(w, r)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}

func initBinaryService(t *testing.T,
	repo map[string]models.BinaryResponse,
) (*services.BinaryService, map[string]models.BinaryResponse) {
	s := services.NewBinaryService(initDataMS(t))
	newRepo := make(map[string]models.BinaryResponse, len(repo))
	for iid, v := range repo {
		id, err := s.StoreBinary(context.Background(), v.UID, models.BinaryRequest{
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

func initDataMS(t *testing.T) data.Service {
	ds, err := data.NewService("")
	if err != nil {
		t.Fatal(err)
	}
	return ds
}
