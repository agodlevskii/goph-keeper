package enc

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrPasswordLength = errors.New("enc: the password is missing")

func HashPassword(s string) (string, error) {
	if len(s) == 0 {
		return "", ErrPasswordLength
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func VerifyPassword(pwd, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd)); err != nil {
		return false
	}
	return true
}
