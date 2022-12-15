package user

import (
	"errors"
)

var ErrDBMissingURL = errors.New("users db url is missing")
var ErrNotFound = errors.New("user not found")

func NewRepo(repoURL string) (IRepository, error) {
	if repoURL == "" {
		return NewBasicRepo(), nil
	}
	return NewDBRepo(repoURL)
}
