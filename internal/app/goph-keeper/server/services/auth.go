package services

import (
	"errors"
)

type AuthReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type AuthService struct {
	ss SessionService
	us UserService
}

func NewAuthService(ss SessionService, us UserService) AuthService {
	return AuthService{
		ss: ss,
		us: us,
	}
}

func (s AuthService) Authorize(token string) (string, error) {
	if exp, err := s.ss.IsTokenExpired(token); err != nil || exp {
		return "", err
	}
	return s.ss.GetUidFromToken(token)
}

func (s AuthService) Login(cid string, u AuthReq) (string, string, error) {
	if cid != "" {
		t, err := s.ss.RestoreSession(cid)
		if err == nil {
			return t, cid, nil
		}
		if err.Error() != "token is expired" && err.Error() != "token not found" {
			return "", "", err
		}
	}

	su, err := s.us.GetUser(u)
	if err != nil {
		if err.Error() == "user not found" {
			return "", "", errors.New("invalid username or password")
		}
		return "", "", err
	}

	token, err := s.ss.GenerateToken(su.ID)
	if err != nil {
		return "", "", err
	}

	cid, err = s.ss.StoreSession(token)
	if err != nil {
		return "", "", err
	}
	return token, cid, nil
}

func (s AuthService) Logout(cid string) (bool, error) {
	if err := s.ss.DeleteSession(cid); err != nil {
		return false, err
	}
	return true, nil
}

func (s AuthService) Register(u AuthReq) error {
	return s.us.AddUser(u)
}
