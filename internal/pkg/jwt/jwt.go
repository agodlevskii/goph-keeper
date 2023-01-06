package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var secret = []byte("j6hsdQ$pj9_ymLQ0")

var (
	ErrTokenExpired = errors.New("jwt: token expired")
	ErrTokenSigning = errors.New("jwt: unexpected token signing method")
	ErrTokenClaims  = errors.New("jwt: failed to extract claims from a token")
)

// EncodeToken creates a token string with encoded user ID and expiry time.
func EncodeToken(uid string, expTime time.Duration) (string, error) {
	if uid == "" {
		return "", ErrTokenClaims
	}
	if expTime == 0 {
		expTime = time.Hour * 12
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(expTime).Unix(),
		"sub": uid,
	})

	return token.SignedString(secret)
}

// GetUserIDFromToken returns the encoded user ID from a token string.
func GetUserIDFromToken(token string) (string, error) {
	claims, err := getClaims(token)
	if err != nil {
		return "", err
	}
	return claims["sub"].(string), nil
}

// IsTokenExpired checks if a token is expired.
func IsTokenExpired(token string) (bool, error) {
	claims, err := getClaims(token)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return true, nil
		}
		return true, err
	}
	return !claims.VerifyExpiresAt(jwt.TimeFunc().Unix(), true), nil
}

func getClaims(ts string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(ts, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenSigning
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return jwt.MapClaims{}, ErrTokenClaims
}
