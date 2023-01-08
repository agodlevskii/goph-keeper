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

// NewService returns an instance of the Service with the associated repository.
func NewService(repoURL string) (Service, error) {
	db, err := NewRepo(repoURL)
	return Service{db: db}, err
}

// RestoreSession gathers the stored client-associated token.
// If the token is expired, the method deletes it from the repository and returns an error.
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

// StoreSession generates new client ID and associates the passed token with it.
func (s Service) StoreSession(ctx context.Context, token string) (string, error) {
	cid := generateClientID()
	return cid, s.db.StoreSession(ctx, cid, token)
}

// DeleteSession deletes the client-associated session.
func (s Service) DeleteSession(ctx context.Context, cid string) error {
	return s.db.DeleteSession(ctx, cid)
}

// GenerateToken generates a new JWT token with the specified expiry time.
func (s Service) GenerateToken(uid string) (string, error) {
	if uid == "" {
		return "", ErrEmptyUID
	}
	return jwt.EncodeToken(uid, 0)
}

// GetUIDFromToken parses the token string and returns the UID from its claims.
func (s Service) GetUIDFromToken(token string) (string, error) {
	if token == "" {
		return "", ErrEmptyToken
	}
	return jwt.GetUserIDFromToken(token)
}

// IsTokenExpired checks if the token had expired.
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
