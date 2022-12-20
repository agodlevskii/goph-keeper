package session

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicRepo_DeleteSession(t *testing.T) {
	for _, tt := range getDeleteSessionCases() {
		t.Run(tt.name, func(t *testing.T) {
			r := initBasicRepo(tt.repo)
			err := r.DeleteSession(context.Background(), tt.cid)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestBasicRepo_GetSession(t *testing.T) {
	for _, tt := range getGetSessionCases() {
		t.Run(tt.name, func(t *testing.T) {
			r := initBasicRepo(tt.repo)
			got, err := r.GetSession(context.Background(), tt.cid)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestBasicRepo_StoreSession(t *testing.T) {
	for _, tt := range getStoreSessionCases() {
		t.Run(tt.name, func(t *testing.T) {
			r := initBasicRepo(tt.repo)
			err := r.StoreSession(context.Background(), tt.args.cid, tt.args.token)
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
			wantField:     "tokens",
			wantFieldType: "*sync.Map",
			wantType:      "*session.BasicRepo",
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
