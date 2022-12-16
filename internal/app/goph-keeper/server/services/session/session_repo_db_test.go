package session

import (
	"context"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDBRepo_DeleteSession(t *testing.T) {
	for _, tt := range getDeleteSessionCases() {
		t.Run(tt.name, func(t *testing.T) {
			r, mock, err := initDBRepo()
			if err != nil {
				t.Fatal(err)
			}

			eq := mock.ExpectExec(regexp.QuoteMeta(DeleteSession)).WithArgs(tt.cid)
			if tt.repo[tt.cid] == "" {
				eq.WillReturnError(ErrNotFound)
			} else {
				eq.WillReturnResult(sqlmock.NewResult(1, 1))
			}

			err = r.DeleteSession(context.Background(), tt.cid)
			assert.Equal(t, tt.wantErr, err)
			checkMetExpectations(t, mock)
		})
	}
}

func TestDBRepo_GetSession(t *testing.T) {
	for _, tt := range getGetSessionCases() {
		t.Run(tt.name, func(t *testing.T) {
			r, mock, err := initDBRepo()
			if err != nil {
				t.Fatal(err)
			}

			eq := mock.ExpectQuery(regexp.QuoteMeta(GetSession)).WithArgs(tt.cid)
			if tt.repo[tt.cid] != "" {
				rows := mock.NewRows([]string{"token"}).AddRow(tt.want)
				eq.WillReturnRows(rows)
			} else {
				eq.WillReturnError(ErrNotFound)
			}

			got, err := r.GetSession(context.Background(), tt.cid)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
			checkMetExpectations(t, mock)
		})
	}
}

func TestDBRepo_StoreSession(t *testing.T) {
	for _, tt := range getStoreSessionCases() {
		t.Run(tt.name, func(t *testing.T) {
			r, mock, err := initDBRepo()
			if err != nil {
				t.Fatal(err)
			}

			_ = getStoreSessionExec(mock, tt.repo, tt.args.cid, tt.args.token)
			err = r.StoreSession(context.Background(), tt.args.cid, tt.args.token)
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
				repoType:  "*session.DBRepo",
				fieldName: "db",
				fieldType: "*sql.DB",
			},
		},
		{
			name: "Wrong Repo URL is present",
			url:  "postgres://localhost:5432/test",
			want: want{
				repoType:  "*session.DBRepo",
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

func getStoreSessionExec(mock sqlmock.Sqlmock, repo map[string]string, cid, token string) *sqlmock.ExpectedExec {
	eq := mock.ExpectExec(regexp.QuoteMeta(StoreSession)).WithArgs(cid, token)
	if cid == "" || token == "" {
		return eq.WillReturnError(ErrIncorrectData)
	}
	if repo[cid] != "" {
		return eq.WillReturnError(ErrSessionExists)
	}
	return eq.WillReturnResult(sqlmock.NewResult(1, 1))
}

func checkMetExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
