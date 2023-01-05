package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/agodlevskii/goph-keeper/internal/pkg/services/session"
	"github.com/agodlevskii/goph-keeper/internal/pkg/services/user"
)

type Service struct {
	sessionService session.Service
	userService    user.Service
}

var (
	ErrSessionExpired  = errors.New("the session has expired, please re-login")
	ErrWrongCredential = errors.New("invalid username or password")
)

func NewService(ss session.Service, us user.Service) Service {
	return Service{
		sessionService: ss,
		userService:    us,
	}
}

func (s Service) Authorize(token string) (string, error) {
	if token == "" {
		return "", ErrWrongCredential
	}

	if exp, err := s.sessionService.IsTokenExpired(token); err != nil || exp {
		if errors.Is(err, session.ErrTokenExpired) || exp {
			return "", ErrSessionExpired
		}
		return "", err
	}
	return s.sessionService.GetUIDFromToken(token)
}

func (s Service) Login(ctx context.Context, cid string, req Payload) (string, string, error) {
	if cid != "" {
		t, err := s.sessionService.RestoreSession(ctx, cid)
		if err == nil {
			return t, cid, nil
		}
		if !errors.Is(err, session.ErrTokenExpired) && !errors.Is(err, session.ErrNotFound) {
			return "", "", err
		}
	}

	u := getUserFromRequest(req)
	su, err := s.userService.GetUser(ctx, u)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return "", "", ErrWrongCredential
		}
		return "", "", err
	}

	token, err := s.sessionService.GenerateToken(su.ID)
	if err != nil {
		return "", "", err
	}

	cid, err = s.sessionService.StoreSession(ctx, token)
	if err != nil {
		return "", "", err
	}
	return token, cid, nil
}

func (s Service) Logout(ctx context.Context, cid string) (bool, error) {
	if cid == "" {
		return false, ErrWrongCredential
	}
	if err := s.sessionService.DeleteSession(ctx, cid); err != nil {
		if errors.Is(err, session.ErrNotFound) {
			return false, ErrWrongCredential
		}
		return false, err
	}
	return true, nil
}

func (s Service) Register(ctx context.Context, req Payload) error {
	u := getUserFromRequest(req)
	return s.userService.AddUser(ctx, u)
}

func getUserFromRequest(req Payload) user.User {
	return user.User{
		Name:     strings.ToLower(req.Name),
		Password: req.Password,
	}
}
