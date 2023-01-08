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

// NewService returns an instance of the Service with the associated repository.
// The repository gets created in accordance with the passed URL.
func NewService(repoURL string) (Service, error) {
	db, err := NewRepo(repoURL)
	return Service{db: db}, err
}

// GetAllDataByType returns all the user's stored data.
func (s Service) GetAllDataByType(ctx context.Context, uid string, t StorageType) ([]SecureData, error) {
	return s.db.GetAllDataByType(ctx, uid, t)
}

// GetDataByID returns the stored data by the unique ID.
// The method returns the data of the specified user only.
func (s Service) GetDataByID(ctx context.Context, uid, id string) (SecureData, error) {
	return s.db.GetDataByID(ctx, uid, id)
}

// StoreSecureDataFromPayload processes payload of any type into a slice of bytes,
// encodes the slice, and stores the content in the DB.
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

// DeleteSecureData removes the stored data with the unique ID.
// The method removes the data of the specified user only.
func (s Service) DeleteSecureData(ctx context.Context, uid, id string) error {
	return s.db.DeleteData(ctx, uid, id)
}

// GetDataFromBytes transforms the encrypted slice of bytes into the original one.
func (s Service) GetDataFromBytes(b []byte) ([]byte, error) {
	return enc.DecryptData(b)
}
