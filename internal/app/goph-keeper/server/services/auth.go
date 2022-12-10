package services

import (
	"errors"

	"github.com/segmentio/ksuid"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/storage"
	"github.com/agodlevskii/goph-keeper/internal/pkg/enc"
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

	su, err := db.GetUserByName(u.Name)
	if err != nil {
		if err.Error() == "user not found" {
			return "", "", errors.New("invalid username or password")
		}
		return "", "", err
	}

	if !enc.VerifyPassword(u.Password, su.Password) {
		return "", "", errors.New("invalid username or password")
	}

	token, err := jwt.EncodeToken(su.ID)
	if err != nil {
		return "", "", err
	}

	cid = generateClientID()
	if err = db.StoreSession(cid, token); err != nil {
		return "", "", err
	}

	return token, cid, nil
}

func Register(db storage.IRepository, u AuthReq) error {
	su, err := db.GetUserByName(u.Name)
	if err != nil && err.Error() != "user not found" {
		return err
	}

	if su.ID != "" {
		return errors.New("user with the specified name already exists")
	}

	hash, err := enc.HashPassword(u.Password)
	if err != nil {
		return nil
	}

	su, err = db.AddUser(u.Name, hash)
	if err != nil {
		return err
	}
	return nil
}

func RestoreSession(db storage.IRepository, cid string) (string, error) {
	t, err := db.GetSession(cid)
	if err != nil {
		return "", err
	}

	if exp, eErr := jwt.IsTokenExpired(t); eErr != nil || exp {
		_ = db.DeleteSession(cid)
		if eErr != nil {
			return "", eErr
		}
		return "", errors.New("token is expired")
	}

	return t, nil
}

func getUserFromRequest(r AuthReq) storage.User {
	return storage.User{
		Name:     r.Name,
		Password: r.Password,
	}
}

func generateClientID() string {
	return ksuid.New().String()
}
