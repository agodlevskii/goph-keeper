package session

import (
	"context"
	"github.com/segmentio/ksuid"

	"github.com/agodlevskii/goph-keeper/internal/pkg/jwt"
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
	if err != nil {
		return Service{}, err
	}
	return Service{db: db}, nil
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
		return "", jwt.ErrTokenExpired
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
	return jwt.EncodeToken(uid)
}

func (s Service) GetUidFromToken(token string) (string, error) {
	return jwt.GetUserIDFromToken(token)
}

func (s Service) IsTokenExpired(token string) (bool, error) {
	if exp, err := jwt.IsTokenExpired(token); err != nil || exp {
		if err != nil {
			return true, err
		}
		return true, jwt.ErrTokenExpired
	}
	return false, nil
}

func generateClientID() string {
	return ksuid.New().String()
}
