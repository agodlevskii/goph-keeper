package text

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

func (s Service) DeleteText(ctx context.Context, uid, id string) error {
	return s.dataService.DeleteSecureData(ctx, uid, id)
}

func (s Service) GetAllTexts(ctx context.Context, uid string) ([]Response, error) {
	sd, err := s.dataService.GetAllDataByType(ctx, uid, data.SText)
	if err != nil {
		if errors.Is(err, data.ErrNotFound) {
			err = ErrNotFound
		}
		return nil, err
	}

	texts := make([]Response, 0, len(sd))
	for _, d := range sd {
		t, dErr := s.getTextFromSecureData(d)
		if dErr != nil {
			return nil, err
		}
		t.Data = ""
		texts = append(texts, t)
	}

	return texts, nil
}

func (s Service) GetTextByID(ctx context.Context, uid, id string) (Response, error) {
	sd, err := s.dataService.GetDataByID(ctx, uid, id)
	if err != nil {
		if errors.Is(err, data.ErrNotFound) {
			err = ErrNotFound
		}
		return Response{}, err
	}
	return s.getTextFromSecureData(sd)
}

func (s Service) StoreText(ctx context.Context, uid string, req Request) (string, error) {
	text := s.getTextFromRequest(uid, req)
	return s.dataService.StoreSecureDataFromPayload(ctx, uid, text, data.SText)
}

func (s Service) getTextFromSecureData(d data.SecureData) (Response, error) {
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

func (s Service) getTextFromRequest(uid string, req Request) Response {
	return Response{
		UID:  uid,
		Name: req.Name,
		Data: req.Data,
		Note: req.Note,
	}
}
