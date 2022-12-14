package client_config

import (
	"fmt"
	"github.com/agodlevskii/goph-keeper/internal/pkg/cfg"
	log "github.com/sirupsen/logrus"
)

type ClientConfig struct {
	File string `env:"CLIENT_CONFIG_FILE" envDefault:"client.yml"`
	API  struct {
		Host   string `json:"host" yaml:"host" env:"API_HOST"`
		Port   int    `json:"port" yaml:"port" env:"API_PORT" envDefault:"8081"`
		Route  string `json:"route" yaml:"route" env:"API_ROUTE" envDefault:"/api/v1"`
		Secure bool   `json:"secure" yaml:"secure" env:"CLIENT_SECURE"`
	} `json:"api" yaml:"api"`
}

func New(opts ...func(*ClientConfig)) *ClientConfig {
	config := &ClientConfig{}
	for _, o := range opts {
		o(config)
	}
	return config
}

func WithEnv() func(*ClientConfig) {
	return func(config *ClientConfig) {
		if err := cfg.UpdateConfigFromEnv(config); err != nil {
			log.Error(err)
		}
	}
}

func WithFile() func(config *ClientConfig) {
	return func(config *ClientConfig) {
		var fCfg ClientConfig
		if err := cfg.UpdateConfigFromFile(config, &fCfg, config.File); err != nil {
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
