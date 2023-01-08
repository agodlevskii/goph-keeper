package binary

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/agodlevskii/goph-keeper/internal/pkg/services/data"
)

func TestNewService(t *testing.T) {
	bds := initBasicDataService(t)
	tests := []struct {
		name string
		want Service
	}{
		{
			name: "Service creation",
			want: Service{dataService: bds},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewService(bds))
		})
	}
}

func TestService_DeleteBinary(t *testing.T) {
	type args struct {
		uid string
		id  string
	}
	tests := []struct {
		name    string
		ds      data.Service
		repo    map[string]Binary
		args    args
		wantErr error
	}{
		{
			name:    "Arguments are empty",
			wantErr: ErrNotFound,
		},
		{
			name:    "ID is empty",
			args:    args{uid: "test"},
			wantErr: ErrNotFound,
		},
		{
			name:    "User ID is empty",
			args:    args{id: "test"},
			wantErr: ErrNotFound,
		},
		{
			name:    "Data is not present",
			repo:    map[string]Binary{"test1": {ID: "test1", UID: "test1"}},
			args:    args{uid: "test", id: "test"},
			wantErr: ErrNotFound,
		},
		{
			name: "Data is present and deleted",
			repo: map[string]Binary{"test": {ID: "test", UID: "test"}},
			args: args{uid: "test", id: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, ids := initService(t, tt.repo)
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

func TestService_GetAllBinaries(t *testing.T) {
	tests := []struct {
		name    string
		uid     string
		repo    map[string]Binary
		want    []Binary
		wantErr error
	}{
		{
			name:    "Missing UID",
			wantErr: ErrNotFound,
		},
		{
			name: "No data",
			uid:  "test1",
			repo: map[string]Binary{"test": {UID: "test", ID: "test"}},
			want: []Binary{},
		},
		{
			name: "Data found",
			uid:  "test",
			repo: map[string]Binary{
				"test":  {UID: "test", Name: "test"},
				"test1": {UID: "test1", Name: "test1"},
			},
			want: []Binary{{UID: "test", Name: "test"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := initService(t, tt.repo)
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

func TestService_GetBinaryByID(t *testing.T) {
	type args struct {
		uid string
		id  string
	}
	tests := []struct {
		name    string
		args    args
		repo    map[string]Binary
		want    Binary
		wantErr error
	}{
		{
			name:    "Missing arguments",
			wantErr: ErrNotFound,
		},
		{
			name:    "Missing UID",
			args:    args{id: "test"},
			wantErr: ErrNotFound,
		},
		{
			name:    "Missing ID",
			args:    args{uid: "test"},
			wantErr: ErrNotFound,
		},
		{
			name:    "No data",
			args:    args{uid: "test", id: "test"},
			repo:    map[string]Binary{"test1": {UID: "test1", ID: "test1"}},
			wantErr: ErrNotFound,
		},
		{
			name: "Data found",
			args: args{uid: "test", id: "test"},
			repo: map[string]Binary{"test": {ID: "test", UID: "test"}},
			want: Binary{ID: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, ids := initService(t, tt.repo)
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

func TestService_getBinaryFromSecureData(t *testing.T) {
	tests := []struct {
		name    string
		d       data.SecureData
		want    Binary
		wantErr error
	}{
		{
			name:    "Invalid data",
			d:       data.SecureData{UID: "test", ID: "test"},
			wantErr: ErrInvalid,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := initService(t, nil)
			got, err := s.getBinaryFromSecureData(tt.d)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func initService(t *testing.T, repo map[string]Binary) (Service, map[string]Binary) {
	s := Service{dataService: initBasicDataService(t)}
	newRepo := make(map[string]Binary, len(repo))
	for iid, v := range repo {
		id, err := s.StoreBinary(context.Background(), v.UID, v)
		if err != nil {
			t.Fatal(err)
		}
		v.ID = id
		newRepo[iid] = v
	}
	return s, newRepo
}

func initBasicDataService(t *testing.T) data.Service {
	ds, err := data.NewService("")
	if err != nil {
		t.Fatal(err)
	}
	return ds
}
