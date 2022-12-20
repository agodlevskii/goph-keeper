package session

import (
	"reflect"
	"sync"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

type deleteSessionCase struct {
	name    string
	repo    map[string]string
	cid     string
	wantErr error
}

type getSessionCase struct {
	name    string
	repo    map[string]string
	cid     string
	want    string
	wantErr error
}

type storeSessionArgs struct {
	cid   string
	token string
}

type storeSessionCase struct {
	name    string
	repo    map[string]string
	args    storeSessionArgs
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
			want: "*session.BasicRepo",
		},
		{
			name:    "Wrong Repo URL is present",
			repoURL: "postgres://localhost:5432/test",
			want:    "*session.DBRepo",
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

func initBasicRepo(data map[string]string) *BasicRepo {
	tokens := &sync.Map{}
	for cid, token := range data {
		tokens.Store(cid, token)
	}
	return &BasicRepo{tokens: tokens}
}

func initDBRepo() (*DBRepo, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	return &DBRepo{db: db}, mock, err
}

func getDeleteSessionCases() []deleteSessionCase {
	return []deleteSessionCase{
		{
			name:    "No client ID passed",
			repo:    map[string]string{"testID": "testToken"},
			wantErr: ErrNotFound,
		},
		{
			name:    "No client ID present",
			repo:    map[string]string{"testID": "testToken"},
			cid:     "testID0",
			wantErr: ErrNotFound,
		},
		{
			name: "Client ID present",
			repo: map[string]string{"testID": "testToken"},
			cid:  "testID",
		},
	}
}

func getGetSessionCases() []getSessionCase {
	return []getSessionCase{
		{
			name:    "No client ID passed",
			repo:    map[string]string{"testID": "testToken"},
			wantErr: ErrNotFound,
		},
		{
			name:    "No client ID present",
			repo:    map[string]string{"testID": "testToken"},
			cid:     "testID0",
			wantErr: ErrNotFound,
		},
		{
			name: "Client ID present",
			repo: map[string]string{"testID": "testToken"},
			cid:  "testID",
			want: "testToken",
		},
	}
}

func getStoreSessionCases() []storeSessionCase {
	return []storeSessionCase{
		{
			name:    "No arguments passed",
			wantErr: ErrIncorrectData,
		},
		{
			name:    "No client ID passed",
			args:    storeSessionArgs{token: "testToken"},
			wantErr: ErrIncorrectData,
		},
		{
			name:    "No token passed",
			args:    storeSessionArgs{cid: "testID"},
			wantErr: ErrIncorrectData,
		},
		{
			name:    "Client ID exists",
			repo:    map[string]string{"testID": "testToken"},
			args:    storeSessionArgs{cid: "testID", token: "testToken"},
			wantErr: ErrSessionExists,
		},
		{
			name: "All arguments are correct",
			repo: map[string]string{"testID": "testToken"},
			args: storeSessionArgs{cid: "testID0", token: "testToken0"},
		},
	}
}
