package binary

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/services/data"
)

type Service struct {
	dataService data.Service
}

var ErrNotFound = errors.New("requested binary data not found")

func NewService(dataService data.Service) Service {
	return Service{dataService: dataService}
}

func (s Service) DeleteBinary(ctx context.Context, uid, id string) error {
	return s.dataService.DeleteSecureData(ctx, uid, id)
}

func (s Service) GetAllBinaries(ctx context.Context, uid string) ([]Response, error) {
	sd, err := s.dataService.GetAllDataByType(ctx, uid, data.SBinary)
	if err != nil {
		if errors.Is(err, data.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	binaries := make([]Response, 0, len(sd))
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

func (s Service) GetBinaryByID(ctx context.Context, uid, id string) (Response, error) {
	d, err := s.dataService.GetDataByID(ctx, uid, id)
	if err != nil {
		if errors.Is(err, data.ErrNotFound) {
			return Response{}, ErrNotFound
		}
		return Response{}, err
	}
	return s.getBinaryFromSecureData(d)
}

func (s Service) StoreBinary(ctx context.Context, uid string, req Request) (string, error) {
	bin := s.getBinaryFromRequest(uid, req)
	return s.dataService.StoreSecureDataFromPayload(ctx, uid, bin, data.SBinary)
}

func (s Service) getBinaryFromSecureData(d data.SecureData) (Response, error) {
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

func (s Service) getBinaryFromRequest(uid string, req Request) Response {
	return Response{
		UID:  uid,
		Name: req.Name,
		Data: req.Data,
		Note: req.Note,
	}
}
