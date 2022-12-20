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

var ErrNotFound = errors.New("requested card data not found")

func NewService(dataService data.Service) Service {
	return Service{dataService: dataService}
}

func (s Service) DeleteCard(ctx context.Context, uid, id string) error {
	return s.dataService.DeleteSecureData(ctx, uid, id)
}

func (s Service) GetAllCards(ctx context.Context, uid string) ([]Card, error) {
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

func (s Service) StoreCard(ctx context.Context, uid string, card Card) (string, error) {
	return s.dataService.StoreSecureDataFromPayload(ctx, uid, card, data.SCard)
}

func (s Service) getCardFromSecureData(d data.SecureData) (Card, error) {
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
