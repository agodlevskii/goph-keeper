package user

import (
	"context"
	"errors"
	"github.com/agodlevskii/goph-keeper/internal/pkg/enc"
	log "github.com/sirupsen/logrus"
	"strings"
)

type IRepository interface {
	AddUser(ctx context.Context, user User) (User, error)
	DeleteUser(ctx context.Context, uid string) error
	GetUserByID(ctx context.Context, uid string) (User, error)
	GetUserByName(ctx context.Context, name string) (User, error)
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

var ErrUserExists = errors.New("the user with specified name already exists")

func (s Service) AddUser(ctx context.Context, user User) error {
	userExist, err := s.doesUserExist(ctx, user)
	if err != nil {
		return err
	}
	if userExist {
		return ErrUserExists
	}

	hash, err := enc.HashPassword(user.Password)
	if err != nil {
		return nil
	}

	user.Password = hash
	_, err = s.db.AddUser(ctx, user)
	return err
}

func (s Service) GetUser(ctx context.Context, user User) (User, error) {
	su, err := s.db.GetUserByName(ctx, strings.ToLower(user.Name))
	if err != nil {
		return User{}, err
	}

	if !enc.VerifyPassword(user.Password, su.Password) {
		log.Error(user.Password, su.Password)
		return User{}, ErrNotFound
	}
	return su, nil
}

func (s Service) doesUserExist(ctx context.Context, user User) (bool, error) {
	su, err := s.db.GetUserByName(ctx, user.Name)
	if err != nil && errors.Is(err, ErrNotFound) {
		return false, err
	}
	return su.ID != "", nil
}
