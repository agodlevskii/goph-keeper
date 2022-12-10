package services

import (
	"errors"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/storage"
	"github.com/agodlevskii/goph-keeper/internal/pkg/jwt"
)

type AuthReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func Authorize(token string) (string, error) {
	if exp, err := jwt.IsTokenExpired(token); err != nil || exp {
		if err != nil {
			return "", err
		}
		return "", errors.New("token is expired")
	}

	return jwt.GetUserIDFromToken(token)
}

func Login(db storage.IRepository, cid string, u AuthReq) (string, string, error) {
	if cid != "" {
		t, err := RestoreSession(db, cid)
		if err == nil {
			return t, cid, nil
		}
		if err.Error() != "token is expired" && err.Error() != "token not found" {
			return "", "", err
		}
	}

	su, err := GetUser(db, u)
	if err != nil {
		if err.Error() == "user not found" {
			return "", "", errors.New("invalid username or password")
		}
		return "", "", err
	}

	token, err := GenerateToken(su.ID)
	if err != nil {
		return "", "", err
	}

	cid, err = StoreSession(db, token)
	if err != nil {
		return "", "", err
	}
	return token, cid, nil
}

func Register(db storage.IRepository, u AuthReq) error {
	return AddUser(db, u)
}
