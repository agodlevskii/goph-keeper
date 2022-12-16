package data

import (
	"context"
	"encoding/json"

	"github.com/agodlevskii/goph-keeper/internal/pkg/enc"
)

type IRepository interface {
	DeleteData(ctx context.Context, uid, id string) error
	GetAllDataByType(ctx context.Context, uid string, t StorageType) ([]SecureData, error)
	GetDataByID(ctx context.Context, uid, id string) (SecureData, error)
	StoreData(ctx context.Context, data SecureData) (string, error)
}

type Service struct {
	db IRepository
}

func NewService(repoURL string) (Service, error) {
	db, err := NewRepo(repoURL)
	if err != nil {
		return Service{}, err
	}
	return Service{db: db}, nil
}

func (s Service) GetAllDataByType(ctx context.Context, uid string, t StorageType) ([]SecureData, error) {
	return s.db.GetAllDataByType(ctx, uid, t)
}

func (s Service) GetDataByID(ctx context.Context, uid, id string) (SecureData, error) {
	return s.db.GetDataByID(ctx, uid, id)
}

func (s Service) StoreSecureDataFromPayload(ctx context.Context, uid string,
	payload any, t StorageType,
) (string, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	encData, err := enc.EncryptData(data)
	if err != nil {
		return "", err
	}

	sd := SecureData{
		UID:  uid,
		Data: encData,
		Type: t,
	}
	return s.db.StoreData(ctx, sd)
}

func (s Service) DeleteSecureData(ctx context.Context, uid, id string) error {
	return s.db.DeleteData(ctx, uid, id)
}

func (s Service) GetDataFromBytes(b []byte) ([]byte, error) {
	return enc.DecryptData(b)
}
