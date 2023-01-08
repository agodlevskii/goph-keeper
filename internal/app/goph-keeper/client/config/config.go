package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/agodlevskii/goph-keeper/internal/pkg/cert"
	"github.com/agodlevskii/goph-keeper/internal/pkg/configs"
)

type ClientConfig struct {
	File string `env:"CLIENT_CONFIG_FILE" envDefault:"client.yml"`
	API  struct {
		Host   string `json:"host" yaml:"host" env:"API_HOST"`
		Port   int    `json:"port" yaml:"port" env:"API_PORT" envDefault:"8081"`
		Route  string `json:"route" yaml:"route" env:"API_ROUTE" envDefault:"/api/v1"`
		Secure bool   `json:"secure" yaml:"secure" env:"CLIENT_SECURE"`
	} `json:"api" yaml:"api"`
	Cert struct {
		CA   string `json:"ca" yaml:"ca" env:"CA_PATH"`
		Cert string `json:"cert" yaml:"cert" env:"CLIENT_CERT_PATH"`
		Key  string `json:"key" yaml:"key" env:"CLIENT_KEY_PATH"`
	} `json:"cert" yaml:"cert"`
}

func New(opts ...func(*ClientConfig)) *ClientConfig {
	cfg := &ClientConfig{}
	for _, o := range opts {
		o(cfg)
	}
	return cfg
}

func WithEnv() func(*ClientConfig) {
	return func(cfg *ClientConfig) {
		if err := configs.UpdateConfigFromEnv(cfg); err != nil {
			log.Error(err)
		}
	}
}

func WithFile() func(cfg *ClientConfig) {
	return func(cfg *ClientConfig) {
		var fCfg ClientConfig
		if err := configs.UpdateConfigFromFile(cfg, &fCfg, cfg.File); err != nil {
			log.Error(err)
		}
	}
}

func (c *ClientConfig) GetAPIAddress() string {
	protocol := "http"
	if c.API.Secure {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s:%d%s", protocol, c.API.Host, c.API.Port, c.API.Route)
}

func (c *ClientConfig) GetCACertPool() (*x509.CertPool, error) {
	return cert.GetCertificatePool(c.Cert.Cert)
}

func (c *ClientConfig) GetCertificate() (tls.Certificate, error) {
	return cert.GetClientCertificate(c.Cert.Cert, c.Cert.Key)
}
