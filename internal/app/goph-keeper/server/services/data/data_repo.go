package data

import (
	"errors"
)

var ErrDBMissingURL = errors.New("data db url is missing")
var ErrNotFound = errors.New("data not found")

func NewRepo(repoURL string) (IRepository, error) {
	if repoURL == "" {
		return NewBasicRepo(), nil
	}
	return NewDBRepo(repoURL)
}
