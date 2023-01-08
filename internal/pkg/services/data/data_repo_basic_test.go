package data

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicRepo_DeleteData(t *testing.T) {
	for _, tt := range getDeleteDataCases() {
		t.Run(tt.name, func(t *testing.T) {
			r := initBasicRepo(tt.repo)
			err := r.DeleteData(context.Background(), tt.args.uid, tt.args.id)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestBasicRepo_GetAllDataByType(t *testing.T) {
	for _, tt := range getGetAllDataByTypeCases() {
		t.Run(tt.name, func(t *testing.T) {
			r := initBasicRepo(tt.repo)
			got, err := r.GetAllDataByType(context.Background(), tt.args.uid, tt.args.t)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestBasicRepo_GetDataByID(t *testing.T) {
	for _, tt := range getGetDataByIDCases() {
		t.Run(tt.name, func(t *testing.T) {
			r := initBasicRepo(tt.repo)
			got, err := r.GetDataByID(context.Background(), tt.args.uid, tt.args.id)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestBasicRepo_StoreData(t *testing.T) {
	for _, tt := range getStoreDataCases() {
		t.Run(tt.name, func(t *testing.T) {
			r := initBasicRepo(tt.repo)
			got, err := r.StoreData(context.Background(), tt.data)
			assert.Equal(t, tt.wantLen, len(got))
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestNewBasicRepo(t *testing.T) {
	tests := []struct {
		name          string
		wantField     string
		wantFieldType string
		wantType      string
	}{
		{
			name:          "Basic repo is created",
			wantField:     "data",
			wantFieldType: "*sync.Map",
			wantType:      "*data.BasicRepo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBasicRepo()
			rGot := reflect.ValueOf(got)
			assert.Equal(t, tt.wantType, rGot.Type().String())

			rField := reflect.Indirect(rGot).Type().Field(0)
			assert.Equal(t, tt.wantField, rField.Name)
			assert.Equal(t, tt.wantFieldType, rField.Type.String())
		})
	}
}
