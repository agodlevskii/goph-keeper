package data

import (
	"context"
	"database/sql"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDBRepo_DeleteData(t *testing.T) {
	for _, tt := range getDeleteDataCases() {
		t.Run(tt.name, func(t *testing.T) {
			r, mock, err := initDBRepo()
			if err != nil {
				t.Fatal(err)
			}

			if tt.args.uid != "" && tt.args.id != "" {
				ee := mock.ExpectExec(regexp.QuoteMeta(DeleteData)).WithArgs(tt.args.uid, tt.args.id)
				var rows int64
				u := tt.repo[tt.args.id]
				if u.ID != "" && u.UID == tt.args.uid {
					rows = 1
				}
				ee.WillReturnResult(sqlmock.NewResult(1, rows))
			}

			err = r.DeleteData(context.Background(), tt.args.uid, tt.args.id)
			assert.Equal(t, tt.wantErr, err)
			checkMetExpectations(t, mock)
		})
	}
}

func TestDBRepo_GetAllDataByType(t *testing.T) {
	for _, tt := range getGetAllDataByTypeCases() {
		t.Run(tt.name, func(t *testing.T) {
			r, mock, err := initDBRepo()
			if err != nil {
				t.Fatal(err)
			}

			if tt.args.uid != "" {
				eq := mock.ExpectQuery(regexp.QuoteMeta(GetAllDataByType)).WithArgs(tt.args.uid, tt.args.t)
				rows := mock.NewRows([]string{"id", "uid", "data", "type"})
				var rowsLen int
				for _, v := range tt.repo {
					if v.UID == tt.args.uid && v.Type == tt.args.t {
						rows.AddRow(v.ID, v.UID, v.Data, v.Type)
						rowsLen++
					}
				}

				if rowsLen > 0 {
					eq.WillReturnRows(rows)
				} else {
					eq.WillReturnError(sql.ErrNoRows)
				}
			}

			got, err := r.GetAllDataByType(context.Background(), tt.args.uid, tt.args.t)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
			checkMetExpectations(t, mock)
		})
	}
}

func TestDBRepo_GetDataByID(t *testing.T) {
	for _, tt := range getGetDataByIDCases() {
		t.Run(tt.name, func(t *testing.T) {
			r, mock, err := initDBRepo()
			if err != nil {
				t.Fatal(err)
			}

			if tt.args.uid != "" && tt.args.id != "" {
				eq := mock.ExpectQuery(regexp.QuoteMeta(GetDataByID)).WithArgs(tt.args.uid, tt.args.id)
				rows := mock.NewRows([]string{"id", "uid", "data", "type"})
				var rowsLen int
				for _, v := range tt.repo {
					if v.UID == tt.args.uid && v.ID == tt.args.id {
						rows.AddRow(v.ID, v.UID, v.Data, v.Type)
						rowsLen++
					}
				}

				if rowsLen > 0 {
					eq.WillReturnRows(rows)
				} else {
					eq.WillReturnError(sql.ErrNoRows)
				}
			}

			got, err := r.GetDataByID(context.Background(), tt.args.uid, tt.args.id)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
			checkMetExpectations(t, mock)
		})
	}
}

func TestDBRepo_StoreData(t *testing.T) {
	for _, tt := range getStoreDataCases() {
		t.Run(tt.name, func(t *testing.T) {
			r, mock, err := initDBRepo()
			if err != nil {
				t.Fatal(err)
			}

			if tt.data.UID != "" && tt.data.Data != nil {
				eq := mock.ExpectQuery(regexp.QuoteMeta(StoreData)).WithArgs(tt.data.UID, tt.data.Data, tt.data.Type)
				rows := mock.NewRows([]string{"id"}).AddRow("123456789012345678901234567890123456")
				eq.WillReturnRows(rows)
			}

			got, err := r.StoreData(context.Background(), tt.data)
			assert.Equal(t, tt.wantLen, len(got))
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
				repoType:  "*data.DBRepo",
				fieldName: "db",
				fieldType: "*sql.DB",
			},
		},
		{
			name: "Wrong Repo URL is present",
			url:  "postgres://localhost:5432/test",
			want: want{
				repoType:  "*data.DBRepo",
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

func checkMetExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
