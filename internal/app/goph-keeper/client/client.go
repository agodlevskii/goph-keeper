package client

import (
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/models"
)

type KeeperClientConfig interface {
	GetAPIAddress() string
}

type KeeperClient interface {
	AuthClient
	BinaryClient
	CardClient
	TextClient
}

type AuthClient interface {
	Login(user, password string) error
	Logout() error
	Register(user, password string) error
}

type BinaryClient interface {
	DeleteBinary(id string) error
	GetAllBinaries() ([]models.BinaryResponse, error)
	GetBinaryByID(id string) (models.BinaryResponse, error)
	StoreBinary(name string, data []byte, note string) (string, error)
}

type CardClient interface {
	DeleteCard(id string) error
	GetAllCards() ([]models.CardResponse, error)
	GetCardByID(id string) (models.CardResponse, error)
	StoreCard(name, number, holder, expDate, cvv, note string) (string, error)
}

type TextClient interface {
	DeleteText(id string) error
	GetAllTexts() ([]models.TextResponse, error)
	GetTextByID(id string) (models.TextResponse, error)
	StoreText(name, data, note string) (string, error)
}

func NewClient() (KeeperClient, error) {
	return NewHTTPClient()
}
