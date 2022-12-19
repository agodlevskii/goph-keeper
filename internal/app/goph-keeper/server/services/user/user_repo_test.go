package user

import (
	"reflect"
	"sync"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

type addUserCase struct {
	name    string
	repo    map[string]User
	user    User
	want    User
	wantErr error
}

type deleteUserCase struct {
	name    string
	repo    map[string]User
	uid     string
	wantErr error
}

type getUserByIDCase struct {
	name    string
	repo    map[string]User
	uid     string
	want    User
	wantErr error
}

type getUserByNameCase struct {
	name    string
	repo    map[string]User
	uName   string
	want    User
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
			want: "*user.BasicRepo",
		},
		{
			name:    "Wrong Repo URL is present",
			repoURL: "postgres://localhost:5432/test",
			want:    "*user.DBRepo",
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

func initBasicRepo(data map[string]User) *BasicRepo {
	users := &sync.Map{}
	for uid, user := range data {
		users.Store(uid, user)
	}
	return &BasicRepo{users: users}
}

func initDBRepo() (*DBRepo, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	return &DBRepo{db: db}, mock, err
}

func getAddUserCases() []addUserCase {
	tu := User{
		ID:       "testID",
		Name:     "test",
		Password: "test",
	}
	return []addUserCase{
		{
			name:    "No user passed",
			wantErr: ErrCredMissing,
		},
		{
			name:    "User with empty name is passed",
			user:    User{Password: "test"},
			wantErr: ErrCredMissing,
		},
		{
			name:    "User with empty password is passed",
			user:    User{Name: "test"},
			wantErr: ErrCredMissing,
		},
		{
			name:    "User ID exists",
			repo:    map[string]User{tu.ID: tu},
			user:    tu,
			wantErr: ErrExists,
		},
		{
			name: "All arguments are correct",
			user: tu,
			want: tu,
		},
	}
}

func getDeleteUserCases() []deleteUserCase {
	return []deleteUserCase{
		{
			name:    "No user ID passed",
			repo:    map[string]User{"testID": {ID: "testID"}},
			wantErr: ErrNotFound,
		},
		{
			name:    "No user ID present",
			repo:    map[string]User{"testID": {ID: "testID"}},
			uid:     "testID0",
			wantErr: ErrNotFound,
		},
		{
			name: "User ID present",
			repo: map[string]User{"testID": {ID: "testID"}},
			uid:  "testID",
		},
	}
}

func getGetUserByIDCases() []getUserByIDCase {
	tu := User{ID: "testID"}
	return []getUserByIDCase{
		{
			name:    "No user ID passed",
			repo:    map[string]User{"testID": tu},
			wantErr: ErrNotFound,
		},
		{
			name:    "No user ID present",
			repo:    map[string]User{"testID": tu},
			uid:     "testID0",
			wantErr: ErrNotFound,
		},
		{
			name: "User ID present",
			repo: map[string]User{"testID": tu},
			uid:  "testID",
			want: tu,
		},
	}
}

func getGetUserByNameCases() []getUserByNameCase {
	tu := User{ID: "testID", Name: "test"}
	return []getUserByNameCase{
		{
			name:    "No user name passed",
			repo:    map[string]User{"test": tu},
			wantErr: ErrNotFound,
		},
		{
			name:    "No user name present",
			repo:    map[string]User{"test": tu},
			uName:   "test0",
			wantErr: ErrNotFound,
		},
		{
			name:  "User name present",
			repo:  map[string]User{"test": tu},
			uName: "test",
			want:  tu,
		},
	}
}
