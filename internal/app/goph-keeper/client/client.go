package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/client/config"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/models"
)

type KeeperClientConfig interface {
	GetAPIAddress() string
	GetCACertPool() (*x509.CertPool, error)
	GetCertificate() (tls.Certificate, error)
}

type KeeperClient interface {
	AuthClient
	BinaryClient
	CardClient
	PasswordClient
	TextClient
}

type AuthClient interface {
	Login(ctx context.Context, user, password string) error
	Logout(ctx context.Context) error
	Register(ctx context.Context, user, password string) error
}

type BinaryClient interface {
	DeleteBinary(ctx context.Context, id string) error
	GetAllBinaries(ctx context.Context) ([]models.BinaryResponse, error)
	GetBinaryByID(ctx context.Context, id string) (models.BinaryResponse, error)
	StoreBinary(ctx context.Context, name string, data []byte, note string) (string, error)
}

type CardClient interface {
	DeleteCard(ctx context.Context, id string) error
	GetAllCards(ctx context.Context) ([]models.CardResponse, error)
	GetCardByID(ctx context.Context, id string) (models.CardResponse, error)
	StoreCard(ctx context.Context, name, number, holder, expDate, cvv, note string) (string, error)
}

type PasswordClient interface {
	DeletePassword(ctx context.Context, id string) error
	GetAllPasswords(ctx context.Context) ([]models.PasswordResponse, error)
	GetPasswordByID(ctx context.Context, id string) (models.PasswordResponse, error)
	StorePassword(ctx context.Context, name, user, password, note string) (string, error)
}

type TextClient interface {
	DeleteText(ctx context.Context, id string) error
	GetAllTexts(ctx context.Context) ([]models.TextResponse, error)
	GetTextByID(ctx context.Context, id string) (models.TextResponse, error)
	StoreText(ctx context.Context, name, data, note string) (string, error)
}

func NewClient(cfg *config.ClientConfig) (KeeperClient, error) {
	return NewHTTPClient(cfg)
}
