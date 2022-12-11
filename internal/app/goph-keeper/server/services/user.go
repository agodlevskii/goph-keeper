package services

import (
	"context"
	"errors"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/storage"
	"github.com/agodlevskii/goph-keeper/internal/pkg/enc"
	log "github.com/sirupsen/logrus"
	"strings"
)

type UserService struct {
	db storage.IUserRepo
}

func NewUserService(db storage.IUserRepo) UserService {
	return UserService{db: db}
}

func (s UserService) AddUser(ctx context.Context, req AuthReq) error {
	user := getUserFromRequest(req)
	userExist, err := doesUserExist(ctx, s.db, user)
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

	user.Password = hash
	_, err = s.db.AddUser(ctx, user)
	return err
}

func (s UserService) GetUser(ctx context.Context, user AuthReq) (storage.User, error) {
	su, err := s.db.GetUserByName(ctx, strings.ToLower(user.Name))
	if err != nil {
		return storage.User{}, err
	}

	if !enc.VerifyPassword(user.Password, su.Password) {
		log.Error(user.Password, su.Password)
		return storage.User{}, errors.New("user not found")
	}
	return su, nil
}

func getUserFromRequest(r AuthReq) storage.User {
	return storage.User{
		Name:     strings.ToLower(r.Name),
		Password: r.Password,
	}
}

func doesUserExist(ctx context.Context, db storage.IUserRepo, user storage.User) (bool, error) {
	su, err := db.GetUserByName(ctx, user.Name)
	if err != nil && err.Error() != "user not found" {
		return false, err
	}
	return su.ID != "", nil
}
