package data

import (
	"errors"
)

var (
	ErrDBMissingURL = errors.New("data db url is missing")
	ErrNotFound     = errors.New("data not found")
	ErrEmpty        = errors.New("data is missing or empty")
	ErrMissingArgs  = errors.New("user id or data type is not specified")
)

func NewRepo(repoURL string) (IRepository, error) {
	if repoURL == "" {
		return NewBasicRepo(), nil
	}
	return NewDBRepo(repoURL)
}
