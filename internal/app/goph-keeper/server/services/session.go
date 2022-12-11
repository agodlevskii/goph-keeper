package services

import (
	"context"
	"errors"

	"github.com/segmentio/ksuid"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/storage"
	"github.com/agodlevskii/goph-keeper/internal/pkg/jwt"
)

type SessionService struct {
	db storage.ISessionRepo
}

func NewSessionService(db storage.ISessionRepo) SessionService {
	return SessionService{db: db}
}

func (s SessionService) RestoreSession(ctx context.Context, cid string) (string, error) {
	t, err := s.db.GetSession(ctx, cid)
	if err != nil {
		return "", err
	}

	if exp, eErr := jwt.IsTokenExpired(t); eErr != nil || exp {
		_ = s.DeleteSession(ctx, cid)
		if eErr != nil {
			return "", eErr
		}
		return "", errors.New("token is expired")
	}

	return t, nil
}

func (s SessionService) StoreSession(ctx context.Context, token string) (string, error) {
	cid := generateClientID()
	return cid, s.db.StoreSession(ctx, cid, token)
}

func (s SessionService) DeleteSession(ctx context.Context, cid string) error {
	return s.db.DeleteSession(ctx, cid)
}

func (s SessionService) GenerateToken(uid string) (string, error) {
	return jwt.EncodeToken(uid)
}

func (s SessionService) GetUidFromToken(token string) (string, error) {
	return jwt.GetUserIDFromToken(token)
}

func (s SessionService) IsTokenExpired(token string) (bool, error) {
	if exp, err := jwt.IsTokenExpired(token); err != nil || exp {
		if err != nil {
			return true, err
		}
		return true, errors.New("token is expired")
	}
	return false, nil
}

func generateClientID() string {
	return ksuid.New().String()
}
