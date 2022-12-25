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

func (s *AuthService) Authorize(token string) (string, error) {
	return s.authMS.Authorize(token)
}

func (s *AuthService) Login(ctx context.Context, cid string, user models.UserRequest) (string, string, error) {
	token, cid, err := s.authMS.Login(ctx, cid, s.getPayloadFromRequest(user))
	if err != nil {
		if errors.Is(err, auth.ErrWrongCredential) {
			return "", "", ErrWrongCredential
		}
		return "", "", err
	}
	return token, cid, err
}

func (s *AuthService) Logout(ctx context.Context, cid string) (bool, error) {
	return s.authMS.Logout(ctx, cid)
}

func (s *AuthService) Register(ctx context.Context, req models.UserRequest) error {
	return s.authMS.Register(ctx, s.getPayloadFromRequest(req))
}

func (s *AuthService) getPayloadFromRequest(req models.UserRequest) auth.Payload {
	return auth.Payload{
		Name:     req.Name,
		Password: req.Password,
	}
}
