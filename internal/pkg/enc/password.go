package enc

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrPasswordLength = errors.New("enc: the password is missing")

// HashPassword encodes an original password string.
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

// VerifyPassword compares an encrypted password with a plain one.
func VerifyPassword(pwd, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err == nil
}
