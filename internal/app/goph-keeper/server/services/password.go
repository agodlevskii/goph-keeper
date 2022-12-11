package services

import (
	"context"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/storage"
)

type PasswordReq struct {
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	Note     string `json:"note"`
}

type PasswordRes struct {
	UID      string `json:"-"`
	ID       string `json:"id"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	Note     string `json:"note"`
}

func GetAllPasswords(ctx context.Context, db storage.IDataRepo, uid string) ([]PasswordRes, error) {
	encPass, err := db.GetAllDataByType(ctx, uid, storage.SPassword)
	if err != nil {
		return nil, err
	}

	ps := make([]PasswordRes, 0, len(encPass))
	for _, ec := range encPass {
		p, eErr := getPasswordFromSecureData(ec)
		if eErr != nil {
			return nil, eErr
		}

		p.Password = "********"
		ps = append(ps, p)
	}
	return ps, nil
}

func GetPasswordByID(ctx context.Context, db storage.IDataRepo, uid, id string) (PasswordRes, error) {
	ep, err := db.GetDataByID(ctx, uid, id)
	if err != nil {
		return PasswordRes{}, nil
	}
	return getPasswordFromSecureData(ep)
}

func StorePassword(ctx context.Context, db storage.IDataRepo, uid string, req PasswordReq) (string, error) {
	pass := getPasswordFromRequest(uid, req)
	return StoreSecureDataFromPayload(ctx, db, uid, pass, storage.SPassword)
}

func getPasswordFromSecureData(d storage.SecureData) (PasswordRes, error) {
	p, err := GetDataFromBytes(d.Data, storage.SPassword)
	if err != nil {
		return PasswordRes{}, err
	}

	pt := p.(PasswordRes)
	pt.ID = d.ID
	return pt, nil
}

func getPasswordFromRequest(uid string, req PasswordReq) PasswordRes {
	return PasswordRes{
		UID:      uid,
		Name:     req.Name,
		User:     req.User,
		Password: req.Password,
		Note:     req.Note,
	}
}
