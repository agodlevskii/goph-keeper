package services

import (
	"errors"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/storage"
	"github.com/agodlevskii/goph-keeper/internal/pkg/enc"
)

func AddUser(db storage.IUserRepository, req AuthReq) error {
	user := getUserFromRequest(req)
	userExist, err := doesUserExist(db, user)
	if err != nil {
		return err
	}
	if userExist {
		return errors.New("user with the specified name already exists")
	}

	hash, err := enc.HashPassword(user.Password)
	if err != nil {
		return nil
	}

	_, err = db.AddUser(user.Name, hash)
	return err
}

func GetUser(db storage.IUserRepository, u AuthReq) (storage.User, error) {
	su, err := db.GetUserByName(u.Name)
	if err != nil {
		return storage.User{}, err
	}

	if !enc.VerifyPassword(u.Password, su.Password) {
		return storage.User{}, errors.New("user not found")
	}
	return su, nil
}

func getUserFromRequest(r AuthReq) storage.User {
	return storage.User{
		Name:     r.Name,
		Password: r.Password,
	}
}

func doesUserExist(db storage.IUserRepository, user storage.User) (bool, error) {
	su, err := db.GetUserByName(user.Name)
	if err != nil && err.Error() != "user not found" {
		return false, err
	}
	return su.ID != "", nil
}
