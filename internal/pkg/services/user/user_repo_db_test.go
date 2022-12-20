package user

import (
	"context"
	"database/sql"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDBRepo_AddUser(t *testing.T) {
	for _, tt := range getAddUserCases() {
		t.Run(tt.name, func(t *testing.T) {
			r, mock, err := initDBRepo()
			if err != nil {
				t.Fatal(err)
			}

			_ = getAddUserQuery(mock, tt.repo, tt.user)
			got, err := r.AddUser(context.Background(), tt.user)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.Password, got.Password)
			assert.Equal(t, tt.wantErr, err)
			checkMetExpectations(t, mock)
		})
	}
}

func TestDBRepo_DeleteUser(t *testing.T) {
	for _, tt := range getDeleteUserCases() {
		t.Run(tt.name, func(t *testing.T) {
			r, mock, err := initDBRepo()
			if err != nil {
				t.Fatal(err)
			}

			ee := mock.ExpectExec(regexp.QuoteMeta(DeleteUser)).WithArgs(tt.uid)
			var rows int64
			if tt.repo[tt.uid].ID != "" {
				rows = 1
			}
			ee.WillReturnResult(sqlmock.NewResult(1, rows))

			err = r.DeleteUser(context.Background(), tt.uid)
			assert.Equal(t, tt.wantErr, err)
			checkMetExpectations(t, mock)
		})
	}
}

func TestDBRepo_GetUserByID(t *testing.T) {
	for _, tt := range getGetUserByIDCases() {
		t.Run(tt.name, func(t *testing.T) {
			r, mock, err := initDBRepo()
			if err != nil {
				t.Fatal(err)
			}

			if tt.uid != "" {
				eq := mock.ExpectQuery(regexp.QuoteMeta(GetUserByID)).WithArgs(tt.uid)
				u := tt.repo[tt.uid]
				if u.ID != "" {
					rows := mock.NewRows([]string{"id", "name", "password"}).AddRow(u.ID, u.Name, u.Password)
					eq.WillReturnRows(rows)
				} else {
					eq.WillReturnError(sql.ErrNoRows)
				}
			}

			got, err := r.GetUserByID(context.Background(), tt.uid)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
			checkMetExpectations(t, mock)
		})
	}
}

func TestDBRepo_GetUserByName(t *testing.T) {
	for _, tt := range getGetUserByNameCases() {
		t.Run(tt.name, func(t *testing.T) {
			r, mock, err := initDBRepo()
			if err != nil {
				t.Fatal(err)
			}

			if tt.uName != "" {
				eq := mock.ExpectQuery(regexp.QuoteMeta(GetUserByName)).WithArgs(tt.uName)
				u := tt.repo[tt.uName]
				if u.Name != "" {
					rows := mock.NewRows([]string{"id", "name", "password"}).AddRow(u.ID, u.Name, u.Password)
					eq.WillReturnRows(rows)
				} else {
					eq.WillReturnError(sql.ErrNoRows)
				}
			}

			got, err := r.GetUserByName(context.Background(), tt.uName)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
			checkMetExpectations(t, mock)
		})
	}
}

func TestNewDBRepo(t *testing.T) {
	type want struct {
		repoType  string
		fieldName string
		fieldType string
	}
	tests := []struct {
		name    string
		url     string
		want    want
		wantErr bool
	}{
		{
			name:    "Empty repo URL",
			wantErr: true,
			want: want{
				repoType:  "*user.DBRepo",
				fieldName: "db",
				fieldType: "*sql.DB",
			},
		},
		{
			name: "Wrong Repo URL is present",
			url:  "postgres://localhost:5432/test",
			want: want{
				repoType:  "*user.DBRepo",
				fieldName: "db",
				fieldType: "*sql.DB",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDBRepo(tt.url)
			assert.Equal(t, tt.wantErr, err != nil)

			rGot := reflect.ValueOf(got)
			assert.Equal(t, tt.want.repoType, rGot.Type().String())

			rField := reflect.Indirect(rGot).Type().Field(0)
			assert.Equal(t, tt.want.fieldName, rField.Name)
			assert.Equal(t, tt.want.fieldType, rField.Type.String())
		})
	}
}

func getAddUserQuery(mock sqlmock.Sqlmock, repo map[string]User, user User) *sqlmock.ExpectedQuery {
	if user.Name == "" || user.Password == "" {
		return nil
	}

	eq := mock.ExpectQuery(regexp.QuoteMeta(AddUser)).WithArgs(user.Name, user.Password)
	if repo[user.ID].Name != "" {
		return eq.WillReturnError(ErrExists)
	}

	rows := mock.NewRows([]string{"id"}).AddRow("id")
	return eq.WillReturnRows(rows)
}

func checkMetExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
