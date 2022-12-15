package password

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/services/data"
)

type Service struct {
	dataService data.Service
}

var ErrNotFound = errors.New("requested password data not found")

func NewService(dataService data.Service) Service {
	return Service{dataService: dataService}
}

func (s Service) DeletePassword(ctx context.Context, uid, id string) error {
	return s.dataService.DeleteSecureData(ctx, uid, id)
}

func (s Service) GetAllPasswords(ctx context.Context, uid string) ([]Response, error) {
	encPass, err := s.dataService.GetAllDataByType(ctx, uid, data.SPassword)
	if err != nil {
		if errors.Is(err, data.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	ps := make([]Response, 0, len(encPass))
	for _, ec := range encPass {
		p, eErr := s.getPasswordFromSecureData(ec)
		if eErr != nil {
			return nil, eErr
		}

		p.Password = "********"
		ps = append(ps, p)
	}
	return ps, nil
}

func (s Service) GetPasswordByID(ctx context.Context, uid, id string) (Response, error) {
	ep, err := s.dataService.GetDataByID(ctx, uid, id)
	if err != nil {
		if errors.Is(err, data.ErrNotFound) {
			return Response{}, ErrNotFound
		}
		return Response{}, nil
	}
	return s.getPasswordFromSecureData(ep)
}

func (s Service) StorePassword(ctx context.Context, uid string, req Request) (string, error) {
	pass := s.getPasswordFromRequest(uid, req)
	return s.dataService.StoreSecureDataFromPayload(ctx, uid, pass, data.SPassword)
}

func (s Service) getPasswordFromSecureData(d data.SecureData) (Response, error) {
	b, err := s.dataService.GetDataFromBytes(d.Data)
	if err != nil {
		return Response{}, err
	}

	var res Response
	if err = json.Unmarshal(b, &res); err != nil {
		return Response{}, err
	}

	res.ID = d.ID
	return res, nil
}

func (s Service) getPasswordFromRequest(uid string, req Request) Response {
	return Response{
		UID:      uid,
		Name:     req.Name,
		User:     req.User,
		Password: req.Password,
		Note:     req.Note,
	}
}
