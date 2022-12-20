package user

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicRepo_AddUser(t *testing.T) {
	for _, tt := range getAddUserCases() {
		t.Run(tt.name, func(t *testing.T) {
			r := initBasicRepo(tt.repo)
			got, err := r.AddUser(context.Background(), tt.user)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.Password, got.Password)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestBasicRepo_DeleteUser(t *testing.T) {
	for _, tt := range getDeleteUserCases() {
		t.Run(tt.name, func(t *testing.T) {
			r := initBasicRepo(tt.repo)
			err := r.DeleteUser(context.Background(), tt.uid)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestBasicRepo_GetUserByID(t *testing.T) {
	for _, tt := range getGetUserByIDCases() {
		t.Run(tt.name, func(t *testing.T) {
			r := initBasicRepo(tt.repo)
			got, err := r.GetUserByID(context.Background(), tt.uid)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestBasicRepo_GetUserByName(t *testing.T) {
	for _, tt := range getGetUserByNameCases() {
		t.Run(tt.name, func(t *testing.T) {
			r := initBasicRepo(tt.repo)
			got, err := r.GetUserByName(context.Background(), tt.uName)
			assert.Equal(t, tt.want, got)
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
			wantField:     "users",
			wantFieldType: "*sync.Map",
			wantType:      "*user.BasicRepo",
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
