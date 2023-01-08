package auth

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/agodlevskii/goph-keeper/internal/pkg/jwt"
	"github.com/agodlevskii/goph-keeper/internal/pkg/services/session"
	"github.com/agodlevskii/goph-keeper/internal/pkg/services/user"
)

func TestNewService(t *testing.T) {
	ss, _ := initSessionService(t, nil)
	us := initUserService(t, nil)
	tests := []struct {
		name string
		want Service
	}{
		{
			name: "Service creation",
			want: Service{sessionService: ss, userService: us},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewService(ss, us))
		})
	}
}

func TestService_Authorize(t *testing.T) {
	token, err := jwt.EncodeToken("test-valid1", 0)
	if err != nil {
		t.Fatal(err)
	}
	expToken, err := jwt.EncodeToken("test-expired2", time.Hour*-1)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		repo    map[string]string
		token   string
		want    string
		wantErr error
	}{
		{
			name:    "No token passed",
			wantErr: ErrWrongCredential,
		},
		{
			name:    "Expired token",
			token:   expToken,
			wantErr: ErrSessionExpired,
		},
		{
			name:  "Valid token",
			token: token,
			want:  "test-valid1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := initService(t, tt.repo, nil)
			got, aErr := s.Authorize(tt.token)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, aErr)
		})
	}
}

func TestService_Login(t *testing.T) {
	token, err := jwt.EncodeToken("test-user", 0)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		cid string
		req Payload
	}
	type repo struct {
		sessions map[string]string
		users    map[string]user.User
	}
	tests := []struct {
		name    string
		repo    repo
		args    args
		want    int
		want1   int
		wantErr error
	}{
		{
			name: "Wrong CID, wrong user",
			repo: repo{
				sessions: map[string]string{token: "id"},
				users:    map[string]user.User{"test": {Name: "test", Password: "test"}},
			},
			args: args{
				cid: "wrong",
				req: Payload{Name: "test", Password: "wrong"},
			},
			wantErr: ErrWrongCredential,
		},
		{
			name: "Wrong CID, right user",
			repo: repo{
				sessions: map[string]string{token: "id"},
				users:    map[string]user.User{"test": {Name: "test", Password: "test"}},
			},
			args: args{
				cid: "wrong",
				req: Payload{Name: "test", Password: "test"},
			},
			want:  165,
			want1: 27,
		},
		{
			name:  "Right CID",
			repo:  repo{sessions: map[string]string{token: "id"}},
			args:  args{cid: "right"},
			want:  129,
			want1: 27,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, r := initService(t, tt.repo.sessions, tt.repo.users)
			if tt.name == "Right CID" {
				tt.args.cid = r[token]
			}

			got, got1, lErr := s.Login(context.Background(), tt.args.cid, tt.args.req)
			assert.Equal(t, tt.want, len(got))
			assert.Equal(t, tt.want1, len(got1))
			assert.Equal(t, tt.wantErr, lErr)
		})
	}
}

func TestService_Logout(t *testing.T) {
	tests := []struct {
		name    string
		cid     string
		want    bool
		wantErr error
	}{
		{
			name:    "Missing CID",
			wantErr: ErrWrongCredential,
		}, {
			name:    "Incorrect CID",
			cid:     "wrong",
			wantErr: ErrWrongCredential,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := initService(t, nil, nil)
			got, err := s.Logout(context.Background(), tt.cid)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestService_Register(t *testing.T) {
	type args struct {
		req  Payload
		user Payload
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "Empty name",
			args:    args{req: Payload{Password: "test"}},
			wantErr: user.ErrCredMissing,
		},
		{
			name:    "Empty password",
			args:    args{req: Payload{Name: "test"}},
			wantErr: user.ErrCredMissing,
		},
		{
			name: "User exists",
			args: args{
				req:  Payload{Name: "test", Password: "test"},
				user: Payload{Name: "test", Password: "test1"},
			},
			wantErr: user.ErrExists,
		},
		{
			name: "User is registered",
			args: args{
				req:  Payload{Name: "test", Password: "test"},
				user: Payload{Name: "test1", Password: "test1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := initService(t, nil, nil)
			if tt.args.user.Name != "" {
				if err := s.Register(context.Background(), tt.args.user); err != nil {
					t.Fatal(err)
				}
			}
			assert.Equal(t, tt.wantErr, s.Register(context.Background(), tt.args.req))
		})
	}
}

func Test_getUserFromRequest(t *testing.T) {
	tests := []struct {
		name string
		req  Payload
		want user.User
	}{
		{
			name: "Missing payload",
		},
		{
			name: "Correct payload",
			req:  Payload{Name: "TEST", Password: "TEST"},
			want: user.User{Name: "test", Password: "TEST"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getUserFromRequest(tt.req)
			assert.Equal(t, tt.want, got)
		})
	}
}

func initService(t *testing.T, sessions map[string]string, users map[string]user.User) (Service, map[string]string) {
	ss, sRepo := initSessionService(t, sessions)
	us := initUserService(t, users)
	return Service{sessionService: ss, userService: us}, sRepo
}

func initSessionService(t *testing.T, sessions map[string]string) (session.Service, map[string]string) {
	s, err := session.NewService("")
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) == 0 {
		return s, nil
	}

	repo := make(map[string]string, len(sessions))
	for v := range sessions {
		cid, sErr := s.StoreSession(context.Background(), v)
		if sErr != nil {
			t.Fatal(err)
		}
		repo[v] = cid
	}

	return s, repo
}

func initUserService(t *testing.T, users map[string]user.User) user.Service {
	s, err := user.NewService("")
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range users {
		if err = s.AddUser(context.Background(), v); err != nil {
			t.Fatal(err)
		}
	}
	return s
}
