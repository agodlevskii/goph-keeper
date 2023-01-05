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

var (
	ErrInvalid  = errors.New("passed password data is invalid")
	ErrNotFound = errors.New("requested password data not found")
)

// NewService returns an instance of the Service with pre-defined data microservice.
func NewService(dataService data.Service) Service {
	return Service{dataService: dataService}
}

// DeletePassword removes the stored data with the unique ID.
// The method removes the data of the specified user only.
func (s Service) DeletePassword(ctx context.Context, uid, id string) error {
	if uid == "" || id == "" {
		return ErrNotFound
	}

	err := s.dataService.DeleteSecureData(ctx, uid, id)
	if errors.Is(err, data.ErrNotFound) {
		return ErrNotFound
	}
	return err
}

// GetAllPasswords returns all the user's stored passwords.
func (s Service) GetAllPasswords(ctx context.Context, uid string) ([]Password, error) {
	if uid == "" {
		return nil, ErrNotFound
	}

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

// GetPasswordByID returns the stored data by the unique ID.
// The method returns the data of the specified user only.
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

// StorePassword stores the original password via the associated data microservice.
func (s Service) StorePassword(ctx context.Context, pass Password) (string, error) {
	return s.dataService.StoreSecureDataFromPayload(ctx, pass.UID, pass, data.SPassword)
}

func (s Service) getPasswordFromSecureData(d data.SecureData) (Password, error) {
	if len(d.Data) == 0 {
		return Password{}, ErrInvalid
	}

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
