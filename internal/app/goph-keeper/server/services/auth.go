package services

import (
	"errors"
)

type AuthReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type AuthService struct {
	session SessionService
	user    UserService
}

func NewAuthService(ss SessionService, us UserService) AuthService {
	return AuthService{
		session: ss,
		user:    us,
	}
}

func (s AuthService) Authorize(token string) (string, error) {
	if exp, err := s.session.IsTokenExpired(token); err != nil || exp {
		return "", err
	}
	return s.session.GetUidFromToken(token)
}

func (s AuthService) Login(cid string, u AuthReq) (string, string, error) {
	if cid != "" {
		t, err := s.session.RestoreSession(cid)
		if err == nil {
			return t, cid, nil
		}
		if err.Error() != "token is expired" && err.Error() != "token not found" {
			return "", "", err
		}
	}

	su, err := s.user.GetUser(u)
	if err != nil {
		if err.Error() == "user not found" {
			return "", "", errors.New("invalid username or password")
		}
		return "", "", err
	}

	token, err := s.session.GenerateToken(su.ID)
	if err != nil {
		return "", "", err
	}

	cid, err = s.session.StoreSession(token)
	if err != nil {
		return "", "", err
	}
	return token, cid, nil
}

func (s AuthService) Logout(cid string) (bool, error) {
	if err := s.session.DeleteSession(cid); err != nil {
		return false, err
	}
	return true, nil
}

func (s AuthService) Register(u AuthReq) error {
	return s.user.AddUser(u)
}
