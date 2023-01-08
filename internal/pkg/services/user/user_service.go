package user

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/agodlevskii/goph-keeper/internal/pkg/enc"
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

// NewService returns an instance of the Service with the associated repository.
func NewService(repoURL string) (Service, error) {
	db, err := NewRepo(repoURL)
	return Service{db: db}, err
}

var ErrCredMissing = errors.New("the user is missing one or more required fields")

// AddUser hashes the passed user's password and stores a new user.
// If the user with the specified name already exists, it returns an error.
func (s Service) AddUser(ctx context.Context, user User) error {
	if user.Name == "" || user.Password == "" {
		return ErrCredMissing
	}
	userExist, err := s.doesUserExist(ctx, user)
	if err != nil {
		return err
	}
	if userExist {
		return ErrExists
	}

	hash, err := enc.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hash
	_, err = s.db.AddUser(ctx, user)
	return err
}

// GetUser gathers the user by its name and compares the passed and the stored passwords.
// If the passwords don't match, or user is not found in the repository, the methods returns an error.
func (s Service) GetUser(ctx context.Context, user User) (User, error) {
	su, err := s.db.GetUserByName(ctx, strings.ToLower(user.Name))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}

	if !enc.VerifyPassword(user.Password, su.Password) {
		return User{}, ErrNotFound
	}
	return su, nil
}

func (s Service) doesUserExist(ctx context.Context, user User) (bool, error) {
	su, err := s.db.GetUserByName(ctx, strings.ToLower(user.Name))
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, ErrNotFound) {
		return false, nil
	}
	return su.ID != "", err
}
