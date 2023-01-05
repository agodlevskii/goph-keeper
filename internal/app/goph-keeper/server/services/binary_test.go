package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/models"
	"github.com/agodlevskii/goph-keeper/internal/pkg/services/binary"
	"github.com/agodlevskii/goph-keeper/internal/pkg/services/data"
)

func TestBinaryService_DeleteBinary(t *testing.T) {
	type args struct {
		uid string
		id  string
	}
	tests := []struct {
		name    string
		repo    map[string]models.BinaryResponse
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
			repo:    map[string]models.BinaryResponse{"t": {ID: "t", UID: "t", Name: "t", Data: []byte("t")}},
			args:    args{uid: "t1", id: "t1"},
			wantErr: ErrBinaryNotFound,
		},
		{
			name: "Data is present and deleted",
			repo: map[string]models.BinaryResponse{"t": {ID: "t", UID: "t", Name: "t", Data: []byte("t")}},
			args: args{uid: "t", id: "t"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, ids := initBinaryService(t, tt.repo)
			if len(ids) > 0 {
				for id, v := range ids {
					if tt.args.id == id {
						tt.args.id = v.ID
					}
				}
			}

			err := s.DeleteBinary(context.Background(), tt.args.uid, tt.args.id)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestBinaryService_GetAllBinaries(t *testing.T) {
	tests := []struct {
		name    string
		uid     string
		repo    map[string]models.BinaryResponse
		want    []models.BinaryResponse
		wantErr error
	}{
		{
			name:    "Missing UID",
			wantErr: ErrBadArguments,
		},
		{
			name: "No data",
			uid:  "t1",
			repo: map[string]models.BinaryResponse{"t": {UID: "t", ID: "t", Name: "t", Data: []byte("t")}},
			want: []models.BinaryResponse{},
		},
		{
			name: "Data found",
			uid:  "t",
			repo: map[string]models.BinaryResponse{
				"t":  {UID: "t", ID: "t", Name: "t", Data: []byte("t")},
				"t1": {UID: "t1", ID: "t1", Name: "t1", Data: []byte("t1")},
			},
			want: []models.BinaryResponse{{UID: "t", Name: "t"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := initBinaryService(t, tt.repo)
			got, err := s.GetAllBinaries(context.Background(), tt.uid)
			if len(got) == 0 {
				assert.Equal(t, tt.want, got)
			} else {
				assert.Equal(t, tt.want[0].Name, got[0].Name)
			}
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestBinaryService_GetBinaryByID(t *testing.T) {
	type args struct {
		uid string
		id  string
	}
	tests := []struct {
		name    string
		args    args
		repo    map[string]models.BinaryResponse
		want    models.BinaryResponse
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
			repo:    map[string]models.BinaryResponse{"t1": {UID: "t1", ID: "t1", Name: "t1", Data: []byte("t1")}},
			wantErr: ErrBinaryNotFound,
		},
		{
			name: "Data found",
			args: args{uid: "t", id: "t"},
			repo: map[string]models.BinaryResponse{"t": {UID: "t", ID: "t", Name: "t", Data: []byte("t")}},
			want: models.BinaryResponse{ID: "t", Name: "t", Data: []byte("t")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, ids := initBinaryService(t, tt.repo)
			if len(ids) > 0 {
				for id, v := range ids {
					if tt.args.id == id {
						tt.args.id = v.ID
						tt.want.ID = v.ID
					}
				}
			}

			got, err := s.GetBinaryByID(context.Background(), tt.args.uid, tt.args.id)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestBinaryService_getModelFromRequest(t *testing.T) {
	type args struct {
		uid string
		req models.BinaryRequest
	}
	tests := []struct {
		name string
		args args
		want binary.Binary
	}{
		{
			name: "Missing UID",
			args: args{req: models.BinaryRequest{Name: "test"}},
			want: binary.Binary{Name: "test"},
		},
		{
			name: "Correct UID",
			args: args{
				uid: "testID",
				req: models.BinaryRequest{Name: "test"},
			},
			want: binary.Binary{UID: "testID", Name: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := initBinaryService(t, nil)
			assert.Equal(t, tt.want, s.getModelFromRequest(tt.args.uid, tt.args.req))
		})
	}
}

func TestHandler_StoreBinary(t *testing.T) {
	type args struct {
		uid  string
		text models.BinaryRequest
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
			args:    args{uid: "test", text: models.BinaryRequest{Data: []byte("test")}},
			wantErr: ErrBadArguments,
		},
		{
			name:    "Empty data",
			args:    args{uid: "test", text: models.BinaryRequest{Name: "test"}},
			wantErr: ErrBadArguments,
		},
		{
			name: "Data saved",
			args: args{uid: "test", text: models.BinaryRequest{Name: "test", Data: []byte("test")}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := initBinaryService(t, nil)
			_, err := s.StoreBinary(context.Background(), tt.args.uid, tt.args.text)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestBinaryService_getResponseFromModel(t *testing.T) {
	tests := []struct {
		name  string
		model binary.Binary
		want  models.BinaryResponse
	}{
		{
			name:  "Correct model",
			model: binary.Binary{Name: "test"},
			want:  models.BinaryResponse{Name: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := initBinaryService(t, nil)
			assert.Equal(t, tt.want, s.getResponseFromModel(tt.model))
		})
	}
}

func TestNewBinaryService(t *testing.T) {
	ds := initDataMS(t)
	tests := []struct {
		name string
		want *BinaryService
	}{
		{
			name: "Service creation",
			want: &BinaryService{binaryMS: binary.NewService(ds)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewBinaryService(ds))
		})
	}
}

func initBinaryService(t *testing.T,
	repo map[string]models.BinaryResponse,
) (*BinaryService, map[string]models.BinaryResponse) {
	s := BinaryService{binaryMS: binary.NewService(initDataMS(t))}
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
	return &s, newRepo
}

func initDataMS(t *testing.T) data.Service {
	ds, err := data.NewService("")
	if err != nil {
		t.Fatal(err)
	}
	return ds
}
