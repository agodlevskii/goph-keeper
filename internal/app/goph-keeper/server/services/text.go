package services

import (
	"context"
	"errors"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/models"
	"github.com/agodlevskii/goph-keeper/internal/pkg/services/data"
	"github.com/agodlevskii/goph-keeper/internal/pkg/services/text"
)

type TextService struct {
	textMS text.Service
}

var ErrTextNotFound = errors.New("requested text data not found")

// NewTextService returns an instance of the BinaryService with pre-defined text microservice.
func NewTextService(dataMS data.Service) *TextService {
	return &TextService{textMS: text.NewService(dataMS)}
}

// DeleteText removes the stored data with the unique ID.
// The method removes the data of the specified user only.
func (s *TextService) DeleteText(ctx context.Context, uid, id string) error {
	if uid == "" || id == "" {
		return ErrBadArguments
	}
	err := s.textMS.DeleteText(ctx, uid, id)
	if errors.Is(err, text.ErrNotFound) {
		return ErrTextNotFound
	}
	return err
}

// GetAllTexts returns all the user's stored texts.
func (s *TextService) GetAllTexts(ctx context.Context, uid string) ([]models.TextResponse, error) {
	if uid == "" {
		return nil, ErrBadArguments
	}
	resp, err := s.textMS.GetAllTexts(ctx, uid)
	if err != nil {
		if errors.Is(err, text.ErrNotFound) {
			return nil, ErrTextNotFound
		}
		return nil, err
	}

	texts := make([]models.TextResponse, 0, len(resp))
	for _, c := range resp {
		texts = append(texts, s.getResponseFromModel(c))
	}
	return texts, nil
}

// GetTextByID returns the stored data by the unique ID.
// The method returns the data of the specified user only.
func (s *TextService) GetTextByID(ctx context.Context, uid, id string) (models.TextResponse, error) {
	if uid == "" || id == "" {
		return models.TextResponse{}, ErrBadArguments
	}
	res, err := s.textMS.GetTextByID(ctx, uid, id)
	if err != nil {
		if errors.Is(err, text.ErrNotFound) {
			return models.TextResponse{}, ErrTextNotFound
		}
		return models.TextResponse{}, err
	}
	return s.getResponseFromModel(res), nil
}

// StoreText stores the original text via the associated data microservice.
func (s *TextService) StoreText(ctx context.Context, uid string, req models.TextRequest) (string, error) {
	if uid == "" || req.Name == "" || req.Data == "" {
		return "", ErrBadArguments
	}
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
