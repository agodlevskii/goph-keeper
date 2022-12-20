package user

import (
	"errors"
)

var (
	ErrDBMissingURL = errors.New("users db url is missing")
	ErrExists       = errors.New("the user with specified name already exists")
	ErrNotFound     = errors.New("user not found")
)

func NewRepo(repoURL string) (IRepository, error) {
	if repoURL == "" {
		return NewBasicRepo(), nil
	}
	return NewDBRepo(repoURL)
}
