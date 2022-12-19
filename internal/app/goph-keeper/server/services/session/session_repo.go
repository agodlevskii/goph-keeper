package session

import (
	"errors"
)

var (
	ErrDBMissingURL  = errors.New("session db url is missing")
	ErrNotFound      = errors.New("session not found")
	ErrIncorrectData = errors.New("client id or token is not specified")
	ErrSessionExists = errors.New("session for specified client id already exists")
)

func NewRepo(repoURL string) (IRepository, error) {
	if repoURL == "" {
		return NewBasicRepo(), nil
	}
	return NewDBRepo(repoURL)
}
