package services

import (
	"context"
	"errors"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/models"
	"github.com/agodlevskii/goph-keeper/internal/pkg/services/binary"
	"github.com/agodlevskii/goph-keeper/internal/pkg/services/data"
)

type BinaryService struct {
	binaryMS binary.Service
}

var ErrBinaryNotFound = errors.New("requested binary data not found")

func NewBinaryService(dataMS data.Service) *BinaryService {
	return &BinaryService{binaryMS: binary.NewService(dataMS)}
}

func (s *BinaryService) DeleteBinary(ctx context.Context, uid, id string) error {
	return s.binaryMS.DeleteBinary(ctx, uid, id)
}

func (s *BinaryService) GetAllBinaries(ctx context.Context, uid string) ([]models.BinaryResponse, error) {
	resp, err := s.binaryMS.GetAllBinaries(ctx, uid)
	if err != nil {
		return nil, err
	}

	binaries := make([]models.BinaryResponse, 0, len(resp))
	for _, c := range resp {
		binaries = append(binaries, s.getResponseFromModel(c))
	}
	return binaries, nil
}

func (s *BinaryService) GetBinaryByID(ctx context.Context, uid, id string) (models.BinaryResponse, error) {
	resp, err := s.binaryMS.GetBinaryByID(ctx, uid, id)
	return s.getResponseFromModel(resp), err
}

func (s *BinaryService) StoreBinary(ctx context.Context, uid string, binary models.BinaryRequest) (string, error) {
	return s.binaryMS.StoreBinary(ctx, uid, s.getModelFromRequest(uid, binary))
}

func (s *BinaryService) getResponseFromModel(model binary.Binary) models.BinaryResponse {
	return models.BinaryResponse{
		UID:  model.UID,
		ID:   model.ID,
		Name: model.Name,
		Data: model.Data,
		Note: model.Note,
	}
}

func (s *BinaryService) getModelFromRequest(uid string, req models.BinaryRequest) binary.Binary {
	return binary.Binary{
		UID:  uid,
		Name: req.Name,
		Data: req.Data,
		Note: req.Note,
	}
}
