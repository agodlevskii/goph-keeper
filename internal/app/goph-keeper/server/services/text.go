package services

import (
	"context"
	"errors"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/models"

	"github.com/agodlevskii/goph-keeper/internal/pkg/services/text"

	"github.com/agodlevskii/goph-keeper/internal/pkg/services/data"
)

type TextService struct {
	textMS text.Service
}

var ErrNotFound = errors.New("requested text data not found")

func NewTextService(dataMS data.Service) *TextService {
	return &TextService{textMS: text.NewService(dataMS)}
}

func (s *TextService) DeleteText(ctx context.Context, uid, id string) error {
	return s.textMS.DeleteText(ctx, uid, id)
}

func (s *TextService) GetAllTexts(ctx context.Context, uid string) ([]models.TextResponse, error) {
	resp, err := s.textMS.GetAllTexts(ctx, uid)
	if err != nil {
		return nil, err
	}

	texts := make([]models.TextResponse, 0, len(resp))
	for _, c := range resp {
		texts = append(texts, s.getResponseFromModel(c))
	}
	return texts, nil
}

func (s *TextService) GetTextByID(ctx context.Context, uid, id string) (models.TextResponse, error) {
	res, err := s.textMS.GetTextByID(ctx, uid, id)
	return s.getResponseFromModel(res), err
}

func (s *TextService) StoreText(ctx context.Context, uid string, req models.TextRequest) (string, error) {
	return s.textMS.StoreText(ctx, s.getModelFromRequest(uid, req))
}

func (s *TextService) getResponseFromModel(model text.Text) models.TextResponse {
	return models.TextResponse{
		UID:  model.UID,
		ID:   model.ID,
		Name: model.Name,
		Data: model.Data,
		Note: model.Note,
	}
}

func (s *TextService) getModelFromRequest(uid string, req models.TextRequest) text.Text {
	return text.Text{
		UID:  uid,
		Name: req.Name,
		Data: req.Data,
		Note: req.Note,
	}
}
