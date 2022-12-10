package services

import (
	"errors"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/storage"
	"github.com/agodlevskii/goph-keeper/internal/pkg/jwt"
	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
)

func RestoreSession(db storage.IRepository, cid string) (string, error) {
	t, err := db.GetSession(cid)
	if err != nil {
		return "", err
	}

	log.Info(t)
	if exp, eErr := jwt.IsTokenExpired(t); eErr != nil || exp {
		log.Error(err)
		_ = DeleteSession(db, cid)
		if eErr != nil {
			return "", eErr
		}
		return "", errors.New("token is expired")
	}

	return t, nil
}

func StoreSession(db storage.IRepository, token string) (string, error) {
	cid := generateClientID()
	return cid, db.StoreSession(cid, token)
}

func DeleteSession(db storage.IRepository, cid string) error {
	return db.DeleteSession(cid)
}

func GenerateToken(uid string) (string, error) {
	return jwt.EncodeToken(uid)
}

func generateClientID() string {
	return ksuid.New().String()
}
