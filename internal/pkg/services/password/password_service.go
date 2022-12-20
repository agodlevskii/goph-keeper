package password

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/agodlevskii/goph-keeper/internal/pkg/services/data"
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

func (s Service) GetAllPasswords(ctx context.Context, uid string) ([]Password, error) {
	encPass, err := s.dataService.GetAllDataByType(ctx, uid, data.SPassword)
	if err != nil {
		if errors.Is(err, data.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	ps := make([]Password, 0, len(encPass))
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

func (s Service) GetPasswordByID(ctx context.Context, uid, id string) (Password, error) {
	ep, err := s.dataService.GetDataByID(ctx, uid, id)
	if err != nil {
		if errors.Is(err, data.ErrNotFound) {
			return Password{}, ErrNotFound
		}
		return Password{}, nil
	}
	return s.getPasswordFromSecureData(ep)
}

func (s Service) StorePassword(ctx context.Context, pass Password) (string, error) {
	return s.dataService.StoreSecureDataFromPayload(ctx, pass.UID, pass, data.SPassword)
}

func (s Service) getPasswordFromSecureData(d data.SecureData) (Password, error) {
	b, err := s.dataService.GetDataFromBytes(d.Data)
	if err != nil {
		return Password{}, err
	}

	var res Password
	if err = json.Unmarshal(b, &res); err != nil {
		return Password{}, err
	}

	res.ID = d.ID
	return res, nil
}
