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

// NewService returns an instance of the Service with the associated session and user microservices.
func NewService(ss session.Service, us user.Service) Service {
	return Service{
		sessionService: ss,
		userService:    us,
	}
}

// Authorize parses the passed token string and returns the user ID associated with it.
// If the token is empty or expired, the method returns an error.
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

// Login establishes the user session based on the client ID and user credential.
// If the client ID is passed, the method looks for the associated stored session.
// If the client ID is empty, or the associated token is not found or expired, the method performs login by credential.
// If the credential doesn't match, or another unknown error has occured, the method returns an error,
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

// Logout clears the stored token, associated with the passed client ID.
// If the client ID is missing, the method returns an error.
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

// Register stores a new user.
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
