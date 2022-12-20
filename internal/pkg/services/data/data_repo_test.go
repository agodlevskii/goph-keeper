package data

import (
	"reflect"
	"sync"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

type deleteDataArgs struct {
	uid string
	id  string
}

type deleteDataCase struct {
	name    string
	repo    map[string]SecureData
	args    deleteDataArgs
	wantErr error
}

type getAllDataByTypeArgs struct {
	uid string
	t   StorageType
}

type getAllDataByTypeCase struct {
	name    string
	repo    map[string]SecureData
	args    getAllDataByTypeArgs
	want    []SecureData
	wantErr error
}

type getDataByIDArgs struct {
	uid string
	id  string
}

type getDataByIDCase struct {
	name    string
	repo    map[string]SecureData
	args    getDataByIDArgs
	want    SecureData
	wantErr error
}

type storeDataCase struct {
	name    string
	repo    map[string]SecureData
	data    SecureData
	wantLen int
	wantErr error
}

func TestNewRepo(t *testing.T) {
	tests := []struct {
		name    string
		repoURL string
		want    string
		wantErr bool
	}{
		{
			name: "Repo URL is missing",
			want: "*data.BasicRepo",
		},
		{
			name:    "Wrong Repo URL is present",
			repoURL: "postgres://localhost:5432/test",
			want:    "*data.DBRepo",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRepo(tt.repoURL)
			assert.Equal(t, tt.wantErr, err != nil)

			rGot := reflect.ValueOf(got)
			assert.Equal(t, tt.want, rGot.Type().String())
		})
	}
}

func initBasicRepo(data map[string]SecureData) *BasicRepo {
	ds := &sync.Map{}
	for id, d := range data {
		if us, ok := ds.Load(d.UID); !ok {
			sd := &sync.Map{}
			sd.Store(id, d)
			ds.Store(d.UID, Storage{user: sd})
		} else {
			us.(Storage).user.Store(id, d)
		}
	}
	return &BasicRepo{data: ds}
}

func initDBRepo() (*DBRepo, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	return &DBRepo{db: db}, mock, err
}

func getDeleteDataCases() []deleteDataCase {
	return []deleteDataCase{
		{
			name:    "No data ID passed",
			repo:    map[string]SecureData{"testID": {ID: "testID"}},
			wantErr: ErrNotFound,
		},
		{
			name:    "No data ID present",
			repo:    map[string]SecureData{"testID": {UID: "testUser", ID: "testID"}},
			args:    deleteDataArgs{uid: "testUser", id: "testID0"},
			wantErr: ErrNotFound,
		},
		{
			name:    "No data for UID present",
			repo:    map[string]SecureData{"testID": {UID: "testUser", ID: "testID"}},
			args:    deleteDataArgs{uid: "testUser0", id: "testID"},
			wantErr: ErrNotFound,
		},
		{
			name: "Both UID and ID present",
			repo: map[string]SecureData{"testID": {UID: "testUser", ID: "testID"}},
			args: deleteDataArgs{uid: "testUser", id: "testID"},
		},
	}
}

func getGetAllDataByTypeCases() []getAllDataByTypeCase {
	tr := map[string]SecureData{
		"testID":  {UID: "testUser", ID: "testID", Type: SCard},
		"testID1": {UID: "testUser", ID: "testID1", Type: SPassword},
	}

	return []getAllDataByTypeCase{
		{
			name:    "No user ID passed",
			repo:    tr,
			args:    getAllDataByTypeArgs{t: SCard},
			wantErr: ErrMissingArgs,
		},
		{
			name: "No type passed",
			repo: tr,
			args: getAllDataByTypeArgs{uid: "testUser"},
		},
		{
			name: "No user ID present",
			repo: tr,
			args: getAllDataByTypeArgs{uid: "testUser1", t: SCard},
		},
		{
			name: "No type present",
			repo: tr,
			args: getAllDataByTypeArgs{uid: "testUser", t: SText},
		},
		{
			name: "All arguments present",
			repo: tr,
			args: getAllDataByTypeArgs{uid: "testUser", t: SCard},
			want: []SecureData{{UID: "testUser", ID: "testID", Type: SCard}},
		},
	}
}

func getGetDataByIDCases() []getDataByIDCase {
	td := SecureData{UID: "testUser", ID: "testID", Data: []byte("test")}
	return []getDataByIDCase{
		{
			name:    "No user ID passed",
			args:    getDataByIDArgs{id: "testID"},
			repo:    map[string]SecureData{td.ID: td},
			wantErr: ErrNotFound,
		},
		{
			name:    "No data ID passed",
			args:    getDataByIDArgs{uid: "testUser"},
			repo:    map[string]SecureData{td.ID: td},
			wantErr: ErrNotFound,
		},
		{
			name:    "No data for user present",
			args:    getDataByIDArgs{uid: "testUser1", id: "testID"},
			repo:    map[string]SecureData{td.ID: td},
			wantErr: ErrNotFound,
		},
		{
			name:    "No data ID present",
			args:    getDataByIDArgs{uid: "testUser", id: "testID1"},
			repo:    map[string]SecureData{"testID0": td},
			wantErr: ErrNotFound,
		},
		{
			name: "Data is present",
			args: getDataByIDArgs{uid: "testUser", id: "testID"},
			repo: map[string]SecureData{td.ID: td},
			want: td,
		},
	}
}

func getStoreDataCases() []storeDataCase {
	td := SecureData{UID: "testUser", ID: "testID", Data: []byte("test")}

	return []storeDataCase{
		{
			name:    "No data passed",
			wantErr: ErrEmpty,
		},
		{
			name:    "User with empty user ID is passed",
			data:    SecureData{Data: []byte("test")},
			wantErr: ErrEmpty,
		},
		{
			name:    "All arguments are correct",
			data:    td,
			wantLen: 36,
		},
	}
}
