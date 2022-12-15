package services

import (
	"context"
	"errors"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/storage"
	"github.com/agodlevskii/goph-keeper/internal/pkg/jwt"
)

type AuthReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type AuthService struct {
	session SessionService
	user    UserService
}

var ErrWrongCredential = errors.New("invalid username or password")

func NewAuthService(ss SessionService, us UserService) AuthService {
	return AuthService{
		session: ss,
		user:    us,
	}
}

func (s AuthService) Authorize(token string) (string, error) {
	if exp, err := s.session.IsTokenExpired(token); err != nil || exp {
		return "", err
	}
	return s.session.GetUidFromToken(token)
}

func (s AuthService) Login(ctx context.Context, cid string, u AuthReq) (string, string, error) {
	if cid != "" {
		t, err := s.session.RestoreSession(ctx, cid)
		if err == nil {
			return t, cid, nil
		}
		if !errors.Is(err, jwt.ErrTokenExpired) && !errors.Is(err, storage.ErrNotFound) {
			return "", "", err
		}
	}

	su, err := s.user.GetUser(ctx, u)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return "", "", ErrWrongCredential
		}
		return "", "", err
	}

	token, err := s.session.GenerateToken(su.ID)
	if err != nil {
		return "", "", err
	}

	cid, err = s.session.StoreSession(ctx, token)
	if err != nil {
		return "", "", err
	}
	return token, cid, nil
}

func (s AuthService) Logout(ctx context.Context, cid string) (bool, error) {
	if err := s.session.DeleteSession(ctx, cid); err != nil {
		return false, err
	}
	return true, nil
}

func (s AuthService) Register(ctx context.Context, u AuthReq) error {
	return s.user.AddUser(ctx, u)
}
