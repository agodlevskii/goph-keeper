package text

import (
	"context"
	"encoding/json"
	"errors"

	data2 "github.com/agodlevskii/goph-keeper/internal/pkg/services/data"
)

type Service struct {
	dataService data2.Service
}

var ErrNotFound = errors.New("requested password data not found")

func NewService(dataService data2.Service) Service {
	return Service{dataService: dataService}
}

func (s Service) DeleteText(ctx context.Context, uid, id string) error {
	return s.dataService.DeleteSecureData(ctx, uid, id)
}

func (s Service) GetAllTexts(ctx context.Context, uid string) ([]Text, error) {
	sd, err := s.dataService.GetAllDataByType(ctx, uid, data2.SText)
	if err != nil {
		if errors.Is(err, data2.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	texts := make([]Text, 0, len(sd))
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

func (s Service) GetTextByID(ctx context.Context, uid, id string) (Text, error) {
	sd, err := s.dataService.GetDataByID(ctx, uid, id)
	if err != nil {
		if errors.Is(err, data2.ErrNotFound) {
			return Text{}, ErrNotFound
		}
		return Text{}, err
	}
	return s.getTextFromSecureData(sd)
}

func (s Service) StoreText(ctx context.Context, text Text) (string, error) {
	return s.dataService.StoreSecureDataFromPayload(ctx, text.UID, text, data2.SText)
}

func (s Service) getTextFromSecureData(d data2.SecureData) (Text, error) {
	b, err := s.dataService.GetDataFromBytes(d.Data)
	if err != nil {
		return Text{}, err
	}

	var res Text
	if err = json.Unmarshal(b, &res); err != nil {
		return Text{}, err
	}

	res.ID = d.ID
	return res, nil
}
