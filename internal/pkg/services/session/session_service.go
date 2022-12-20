package session

import (
	"context"
	"errors"

	"github.com/segmentio/ksuid"

	"github.com/agodlevskii/goph-keeper/internal/pkg/jwt"
)

var (
	ErrEmptyToken   = errors.New("session: token is missing or empty")
	ErrEmptyUID     = errors.New("session: uid is missing or empty")
	ErrTokenExpired = errors.New("session: token expired")
)

type IRepository interface {
	DeleteSession(ctx context.Context, cid string) error
	GetSession(ctx context.Context, cid string) (string, error)
	StoreSession(ctx context.Context, cid, token string) error
}

type Service struct {
	db IRepository
}

func NewService(repoURL string) (Service, error) {
	db, err := NewRepo(repoURL)
	return Service{db: db}, err
}

func (s Service) RestoreSession(ctx context.Context, cid string) (string, error) {
	t, err := s.db.GetSession(ctx, cid)
	if err != nil {
		return "", err
	}

	if exp, eErr := jwt.IsTokenExpired(t); eErr != nil || exp {
		_ = s.DeleteSession(ctx, cid)
		if eErr != nil {
			return "", eErr
		}
		return "", ErrTokenExpired
	}

	return t, nil
}

func (s Service) StoreSession(ctx context.Context, token string) (string, error) {
	cid := generateClientID()
	return cid, s.db.StoreSession(ctx, cid, token)
}

func (s Service) DeleteSession(ctx context.Context, cid string) error {
	return s.db.DeleteSession(ctx, cid)
}

func (s Service) GenerateToken(uid string) (string, error) {
	if uid == "" {
		return "", ErrEmptyUID
	}
	return jwt.EncodeToken(uid, 0)
}

func (s Service) GetUIDFromToken(token string) (string, error) {
	if token == "" {
		return "", ErrEmptyToken
	}
	return jwt.GetUserIDFromToken(token)
}

func (s Service) IsTokenExpired(token string) (bool, error) {
	if exp, err := jwt.IsTokenExpired(token); err != nil || exp {
		if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
			return true, err
		}
		return true, ErrTokenExpired
	}
	return false, nil
}

func generateClientID() string {
	return ksuid.New().String()
}
