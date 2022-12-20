package client

import (
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/models"
)

type KeeperClientConfig interface {
	GetAPIAddress() string
}

type KeeperClient interface {
	Login(user, password string) error
	Logout() error
	Register(user, password string) error

	DeleteBinary(id string) error
	GetAllBinaries() ([]models.BinaryResponse, error)
	GetBinaryByID(id string) (models.BinaryResponse, error)
	StoreBinary(name string, data []byte, note string) (string, error)
}

func NewClient() (KeeperClient, error) {
	return NewHTTPClient()
}
