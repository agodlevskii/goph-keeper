package card

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
	ErrInvalid  = errors.New("passed card data is invalid")
	ErrNotFound = errors.New("requested card data not found")
)

// NewService returns an instance of the Service with pre-defined data microservice.
func NewService(dataService data.Service) Service {
	return Service{dataService: dataService}
}

// DeleteCard removes the stored data with the unique ID.
// The method removes the data of the specified user only.
func (s Service) DeleteCard(ctx context.Context, uid, id string) error {
	if uid == "" || id == "" {
		return ErrNotFound
	}

	err := s.dataService.DeleteSecureData(ctx, uid, id)
	if errors.Is(err, data.ErrNotFound) {
		return ErrNotFound
	}
	return err
}

// GetAllCards returns all the user's stored cards.
func (s Service) GetAllCards(ctx context.Context, uid string) ([]Card, error) {
	if uid == "" {
		return nil, ErrNotFound
	}

	sd, err := s.dataService.GetAllDataByType(ctx, uid, data.SCard)
	if err != nil {
		if errors.Is(err, data.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	cards := make([]Card, 0, len(sd))
	for _, d := range sd {
		c, eErr := s.getCardFromSecureData(d)
		if eErr != nil {
			return nil, eErr
		}

		c.CVV = "***"
		cards = append(cards, c)
	}
	return cards, nil
}

// GetCardByID returns the stored data by the unique ID.
// The method returns the data of the specified user only.
func (s Service) GetCardByID(ctx context.Context, uid, id string) (Card, error) {
	d, err := s.dataService.GetDataByID(ctx, uid, id)
	if err != nil {
		if errors.Is(err, data.ErrNotFound) {
			return Card{}, ErrNotFound
		}
		return Card{}, nil
	}
	return s.getCardFromSecureData(d)
}

// StoreCard stores the original card via the associated data microservice.
func (s Service) StoreCard(ctx context.Context, card Card) (string, error) {
	return s.dataService.StoreSecureDataFromPayload(ctx, card.UID, card, data.SCard)
}

func (s Service) getCardFromSecureData(d data.SecureData) (Card, error) {
	if len(d.Data) == 0 {
		return Card{}, ErrInvalid
	}

	b, err := s.dataService.GetDataFromBytes(d.Data)
	if err != nil {
		return Card{}, err
	}

	var res Card
	if err = json.Unmarshal(b, &res); err != nil {
		return Card{}, err
	}

	res.ID = d.ID
	return res, nil
}
