package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/services"
)

type httpRes struct {
	code        int
	resp        string
	contentType string
}

const (
	userCookieName   = "uid"
	clientCookieName = "cid"
	authURL          = "/api/v1/auth"
	binaryURL        = "/api/v1/storage/binary"
	cardURL          = "/api/v1/storage/card"
	pStorageURL      = "/api/v1/storage/password"
	textURL          = "/api/v1/storage/text"
)

func initTestRequest(t *testing.T, method, url, id, uid string, data any) *http.Request {
	rctx := chi.NewRouteContext()
	if id != "" {
		url = fmt.Sprintf("%s/%s", url, id)
		rctx.URLParams.Add("id", id)
	}

	var bReader io.Reader
	if data != nil {
		body, err := json.Marshal(data)
		if err != nil {
			t.Fatal(err)
		}
		if len(body) > 0 {
			bReader = bytes.NewBuffer(body)
		}
	}

	r := httptest.NewRequest(method, url, bReader)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	r = r.WithContext(context.WithValue(r.Context(), uidKey, uid))
	return r
}

func TestNewHandler(t *testing.T) {
	type want struct {
		pattern  string
		handlers int
	}
	tests := []struct {
		name    string
		repoURL string
		want    want
		wantErr bool
	}{
		{
			name: "Init handler",
			want: want{
				pattern:  "/api/v1/*",
				handlers: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHandler(tt.repoURL)
			log.Info(len(got.Routes()[0].Handlers))
			if len(got.Routes()) > 0 {
				route := got.Routes()[0]
				assert.Equal(t, tt.want.pattern, route.Pattern)
				assert.Equal(t, tt.want.handlers, len(route.Handlers))
			}
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func Test_initHandler(t *testing.T) {
	type want struct {
		authType string
	}
	tests := []struct {
		name    string
		repoURL string
		want    want
		wantErr bool
	}{
		{
			name: "Init handler",
			want: want{
				authType: "handlers.IAuthService",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := initHandler(tt.repoURL)
			assert.Equal(t, tt.wantErr, err != nil)

			if err == nil {
				rGot := reflect.ValueOf(got)
				rField := reflect.Indirect(rGot).Type().Field(0)
				assert.Equal(t, tt.want.authType, rField.Type.String())
			}
		})
	}
}

func TestHandler_getErrorCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{
			name: "No error",
			want: http.StatusInternalServerError,
		},
		{
			name: "Unknown error",
			err:  http.ErrAbortHandler,
			want: http.StatusInternalServerError,
		},
		{
			name: "Data not found",
			err:  services.ErrBinaryNotFound,
			want: http.StatusNotFound,
		},
		{
			name: "Bad arguments",
			err:  services.ErrBadArguments,
			want: http.StatusBadRequest,
		},
		{
			name: "Bad credential",
			err:  services.ErrWrongCredential,
			want: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Handler{}
			assert.Equal(t, tt.want, h.getErrorCode(tt.err))
		})
	}
}
