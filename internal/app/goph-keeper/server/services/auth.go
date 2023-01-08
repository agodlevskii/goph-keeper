package services

import (
	"context"
	"errors"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/models"
	"github.com/agodlevskii/goph-keeper/internal/pkg/services/auth"
	"github.com/agodlevskii/goph-keeper/internal/pkg/services/session"
	"github.com/agodlevskii/goph-keeper/internal/pkg/services/user"
)

type AuthService struct {
	authMS auth.Service
}

var ErrWrongCredential = errors.New("invalid username or password")

// NewAuthService returns an instance of the BinaryService with pre-defined auth microservice.
func NewAuthService(repoURL string) (*AuthService, error) {
	sessionMS, err := session.NewService(repoURL)
	if err != nil {
		return nil, err
	}

	userMS, err := user.NewService(repoURL)
	if err != nil {
		return nil, err
	}

	return &AuthService{authMS: auth.NewService(sessionMS, userMS)}, nil
}

// Authorize parses the passed token string and returns the user ID associated with it.
// If the token is empty or expired, the method returns an error.
func (s *AuthService) Authorize(token string) (string, error) {
	if token == "" {
		return "", ErrBadArguments
	}
	return s.authMS.Authorize(token)
}

// Login establishes the user session based on the client ID and user credential.
// If the client ID is passed, the method looks for the associated stored session.
// If the client ID is empty, or the associated token is not found or expired, the method performs login by credential.
// If the credential doesn't match, or another unknown error has occurred, the method returns an error.
func (s *AuthService) Login(ctx context.Context, cid string, user models.UserRequest) (string, string, error) {
	if user.Name == "" || user.Password == "" {
		return "", "", ErrBadArguments
	}
	token, cid, err := s.authMS.Login(ctx, cid, s.getPayloadFromRequest(user))
	if err != nil {
		if errors.Is(err, auth.ErrWrongCredential) {
			return "", "", ErrWrongCredential
		}
		return "", "", err
	}
	return token, cid, err
}

// Logout clears the stored token, associated with the passed client ID.
// If the client ID is missing, the method returns an error.
func (s *AuthService) Logout(ctx context.Context, cid string) (bool, error) {
	ok, err := s.authMS.Logout(ctx, cid)
	if errors.Is(err, auth.ErrWrongCredential) {
		return false, ErrWrongCredential
	}
	return ok, err
}

// Register stores a new user.
func (s *AuthService) Register(ctx context.Context, req models.UserRequest) error {
	if req.Name == "" || req.Password == "" {
		return ErrBadArguments
	}
	return s.authMS.Register(ctx, s.getPayloadFromRequest(req))
}

func (s *AuthService) getPayloadFromRequest(req models.UserRequest) auth.Payload {
	return auth.Payload{
		Name:     req.Name,
		Password: req.Password,
	}
}
