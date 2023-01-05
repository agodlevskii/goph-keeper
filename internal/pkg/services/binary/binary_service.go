package binary

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
	ErrInvalid  = errors.New("passed text data is invalid")
	ErrNotFound = errors.New("requested binary data not found")
)

// NewService returns an instance of the Service with pre-defined data microservice.
func NewService(dataService data.Service) Service {
	return Service{dataService: dataService}
}

// DeleteBinary removes the stored data with the unique ID.
// The method removes the data of the specified user only.
func (s Service) DeleteBinary(ctx context.Context, uid, id string) error {
	err := s.dataService.DeleteSecureData(ctx, uid, id)
	if errors.Is(err, data.ErrNotFound) {
		return ErrNotFound
	}
	return err
}

// GetAllBinaries returns all the user's stored binaries.
func (s Service) GetAllBinaries(ctx context.Context, uid string) ([]Binary, error) {
	sd, err := s.dataService.GetAllDataByType(ctx, uid, data.SBinary)
	if err != nil {
		if errors.Is(err, data.ErrNotFound) || errors.Is(err, data.ErrMissingArgs) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	binaries := make([]Binary, 0, len(sd))
	for _, d := range sd {
		b, dErr := s.getBinaryFromSecureData(d)
		if dErr != nil {
			return nil, err
		}
		b.Data = make([]byte, 0)
		binaries = append(binaries, b)
	}

	return binaries, nil
}

// GetBinaryByID returns the stored data by the unique ID.
// The method returns the data of the specified user only.
func (s Service) GetBinaryByID(ctx context.Context, uid, id string) (Binary, error) {
	if uid == "" || id == "" {
		return Binary{}, ErrNotFound
	}

	d, err := s.dataService.GetDataByID(ctx, uid, id)
	if err != nil {
		if errors.Is(err, data.ErrNotFound) {
			return Binary{}, ErrNotFound
		}
		return Binary{}, err
	}
	return s.getBinaryFromSecureData(d)
}

// StoreBinary stores the original binary via the associated data microservice.
func (s Service) StoreBinary(ctx context.Context, uid string, binary Binary) (string, error) {
	return s.dataService.StoreSecureDataFromPayload(ctx, uid, binary, data.SBinary)
}

func (s Service) getBinaryFromSecureData(d data.SecureData) (Binary, error) {
	if len(d.Data) == 0 {
		return Binary{}, ErrInvalid
	}

	b, err := s.dataService.GetDataFromBytes(d.Data)
	if err != nil {
		return Binary{}, err
	}

	var res Binary
	if err = json.Unmarshal(b, &res); err != nil {
		return Binary{}, err
	}

	res.ID = d.ID
	return res, nil
}
