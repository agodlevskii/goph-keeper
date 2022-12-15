package card

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/services/data"
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

func (s Service) GetAllCards(ctx context.Context, uid string) ([]Response, error) {
	sd, err := s.dataService.GetAllDataByType(ctx, uid, data.SCard)
	if err != nil {
		if errors.Is(err, data.ErrNotFound) {
			err = ErrNotFound
		}
		return nil, err
	}

	cards := make([]Response, 0, len(sd))
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

func (s Service) GetCardByID(ctx context.Context, uid, id string) (Response, error) {
	d, err := s.dataService.GetDataByID(ctx, uid, id)
	if err != nil {
		if errors.Is(err, data.ErrNotFound) {
			err = ErrNotFound
		}
		return Response{}, nil
	}
	return s.getCardFromSecureData(d)
}

func (s Service) StoreCard(ctx context.Context, uid string, req Request) (string, error) {
	card := s.getCardFromRequest(uid, req)
	return s.dataService.StoreSecureDataFromPayload(ctx, uid, card, data.SCard)
}

func (s Service) getCardFromSecureData(d data.SecureData) (Response, error) {
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

func (s Service) getCardFromRequest(uid string, req Request) Response {
	return Response{
		UID:     uid,
		Name:    req.Name,
		Number:  req.Number,
		Holder:  req.Holder,
		ExpDate: req.ExpDate,
		CVV:     req.CVV,
		Note:    req.Note,
	}
}
