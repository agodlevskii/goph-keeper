package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/models"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/services"
	"github.com/agodlevskii/goph-keeper/internal/pkg/jwt"
)

func TestHandler_Auth(t *testing.T) {
	token, err := jwt.EncodeToken("test_id", 0)
	if err != nil {
		t.Fatal(err)
	}
	expToken, err := jwt.EncodeToken("bad_id", -1*time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		cookie *http.Cookie
		want   httpRes
	}{
		{
			name: "Missing cookie",
			want: httpRes{code: http.StatusUnauthorized},
		},
		{
			name:   "Valid token cookie",
			cookie: &http.Cookie{Name: userCookieName, Value: token, Path: "/"},
			want:   httpRes{code: http.StatusOK},
		},
		{
			name:   "Expired token cookie",
			cookie: &http.Cookie{Name: userCookieName, Value: expToken, Path: "/"},
			want:   httpRes{code: http.StatusUnauthorized},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs, _ := initBinaryService(t, nil)
			as, _ := initAuthService(t, models.UserRequest{})
			h := Handler{
				authService:   as,
				binaryService: bs,
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, binaryURL, nil)
			if tt.cookie != nil {
				r.AddCookie(tt.cookie)
			}

			h.Auth(h.GetAllBinaries()).ServeHTTP(w, r)
			got := w.Result()
			assert.Equal(t, tt.want.code, got.StatusCode)
		})
	}
}

func TestHandler_Login(t *testing.T) {
	type fields struct {
		cookie *http.Cookie
		req    models.UserRequest
		user   models.UserRequest
	}
	tests := []struct {
		name   string
		fields fields
		want   httpRes
	}{
		{
			name: "No payload",
			want: httpRes{code: http.StatusBadRequest},
		},
		{
			name:   "No username",
			fields: fields{req: models.UserRequest{Password: "test"}},
			want:   httpRes{code: http.StatusBadRequest},
		},
		{
			name:   "No password",
			fields: fields{req: models.UserRequest{Name: "test"}},
			want:   httpRes{code: http.StatusBadRequest},
		},
		{
			name:   "User doesn't exist",
			fields: fields{req: models.UserRequest{Name: "test", Password: "test"}},
			want:   httpRes{code: http.StatusUnauthorized},
		},
		{
			name: "User exists",
			fields: fields{
				req:  models.UserRequest{Name: "test", Password: "test"},
				user: models.UserRequest{Name: "test", Password: "test"},
			},
			want: httpRes{code: http.StatusOK},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as, _ := initAuthService(t, tt.fields.user)

			h := Handler{authService: as}
			w := httptest.NewRecorder()
			url := fmt.Sprintf("%s/%s", authURL, "register")
			r := initTestRequest(t, http.MethodPost, url, "", "", tt.fields.req)

			h.Login()(w, r)
			got := w.Result()
			assert.Equal(t, tt.want.code, got.StatusCode)
		})
	}
}

func TestHandler_Logout(t *testing.T) {
	type fields struct {
		cookie *http.Cookie
		user   models.UserRequest
	}
	tests := []struct {
		name   string
		fields fields
		want   httpRes
	}{
		{
			name: "Missing cookie",
			want: httpRes{code: http.StatusUnauthorized},
		},
		{
			name: "Valid client cookie",
			fields: fields{
				cookie: &http.Cookie{Name: clientCookieName, Value: "test", Path: "/"},
				user:   models.UserRequest{Name: "test", Password: "test"},
			},
			want: httpRes{code: http.StatusOK},
		},
		{
			name: "Invalid client cookie",
			fields: fields{
				cookie: &http.Cookie{Name: clientCookieName, Value: "test1", Path: "/"},
				user:   models.UserRequest{Name: "test", Password: "test"},
			},
			want: httpRes{code: http.StatusUnauthorized},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as, cid := initAuthService(t, tt.fields.user)
			if cid != "" && tt.name != "Invalid client cookie" {
				tt.fields.cookie.Value = cid
			}

			h := Handler{authService: as}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", authURL, "logout"), nil)
			if tt.fields.cookie != nil {
				r.AddCookie(tt.fields.cookie)
			}

			h.Logout()(w, r)
			got := w.Result()
			assert.Equal(t, tt.want.code, got.StatusCode)
		})
	}
}

func TestHandler_Register(t *testing.T) {
	type fields struct {
		req  models.UserRequest
		user models.UserRequest
	}
	tests := []struct {
		name   string
		fields fields
		want   httpRes
	}{
		{
			name: "No payload",
			want: httpRes{code: http.StatusBadRequest},
		},
		{
			name:   "No username",
			fields: fields{req: models.UserRequest{Password: "test"}},
			want:   httpRes{code: http.StatusBadRequest},
		},
		{
			name:   "No password",
			fields: fields{req: models.UserRequest{Name: "test"}},
			want:   httpRes{code: http.StatusBadRequest},
		},
		{
			name: "User exists",
			fields: fields{
				req:  models.UserRequest{Name: "test", Password: "test"},
				user: models.UserRequest{Name: "test", Password: "test"},
			},
			want: httpRes{code: http.StatusConflict},
		},

		{
			name: "User doesn't exist",
			fields: fields{
				req:  models.UserRequest{Name: "test", Password: "test"},
				user: models.UserRequest{Name: "test1", Password: "test1"},
			},
			want: httpRes{code: http.StatusOK},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as, _ := initAuthService(t, tt.fields.user)

			h := Handler{authService: as}
			w := httptest.NewRecorder()
			url := fmt.Sprintf("%s/%s", authURL, "register")
			r := initTestRequest(t, http.MethodPost, url, "", "", tt.fields.req)

			h.Register()(w, r)
			got := w.Result()
			assert.Equal(t, tt.want.code, got.StatusCode)
		})
	}
}

func initAuthService(t *testing.T, req models.UserRequest) (*services.AuthService, string) {
	as, err := services.NewAuthService("")
	if err != nil {
		t.Fatal(err)
	}

	var cid string
	if req.Name != "" {
		if err = as.Register(context.Background(), req); err != nil {
			t.Fatal(err)
		}
		_, cid, err = as.Login(context.Background(), "", req)
		if err != nil {
			t.Fatal(err)
		}
	}
	return as, cid
}
