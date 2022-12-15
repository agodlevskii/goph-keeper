package session

import (
	"errors"
)

var (
	ErrDBMissingURL = errors.New("session db url is missing")
	ErrNotFound     = errors.New("session not found")
)

func NewRepo(repoURL string) (IRepository, error) {
	if repoURL == "" {
		return NewBasicRepo(), nil
	}
	return NewDBRepo(repoURL)
}
