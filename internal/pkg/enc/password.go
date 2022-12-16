package enc

import "golang.org/x/crypto/bcrypt"

func HashPassword(s string) (string, error) {
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
