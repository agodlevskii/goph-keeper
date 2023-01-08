package text

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
	ErrNotFound = errors.New("requested text data not found")
)

// NewService returns an instance of the Service with pre-defined data microservice.
func NewService(dataService data.Service) Service {
	return Service{dataService: dataService}
}

// DeleteText removes the stored data with the unique ID.
// The method removes the data of the specified user only.
func (s Service) DeleteText(ctx context.Context, uid, id string) error {
	if err := s.dataService.DeleteSecureData(ctx, uid, id); err != nil {
		if errors.Is(err, data.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

// GetAllTexts returns all the user's stored texts.
func (s Service) GetAllTexts(ctx context.Context, uid string) ([]Text, error) {
	if uid == "" {
		return nil, ErrNotFound
	}

	sd, err := s.dataService.GetAllDataByType(ctx, uid, data.SText)
	if err != nil {
		if errors.Is(err, data.ErrNotFound) {
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

// GetTextByID returns the stored data by the unique ID.
// The method returns the data of the specified user only.
func (s Service) GetTextByID(ctx context.Context, uid, id string) (Text, error) {
	if uid == "" || id == "" {
		return Text{}, ErrNotFound
	}

	sd, err := s.dataService.GetDataByID(ctx, uid, id)
	if err != nil {
		if errors.Is(err, data.ErrNotFound) {
			return Text{}, ErrNotFound
		}
		return Text{}, err
	}
	return s.getTextFromSecureData(sd)
}

// StoreText stores the original text via the associated data microservice.
func (s Service) StoreText(ctx context.Context, text Text) (string, error) {
	return s.dataService.StoreSecureDataFromPayload(ctx, text.UID, text, data.SText)
}

func (s Service) getTextFromSecureData(d data.SecureData) (Text, error) {
	if len(d.Data) == 0 {
		return Text{}, ErrInvalid
	}
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
