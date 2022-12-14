package server_config

import (
	"fmt"
	"github.com/agodlevskii/goph-keeper/internal/pkg/cfg"
	log "github.com/sirupsen/logrus"
)

type ServerConfig struct {
	File   string `env:"SERVER_CONFIG_FILE" envDefault:"server.yml"`
	Server struct {
		Host   string `json:"host" yaml:"host" env:"SERVER_HOST"`
		Port   int    `json:"port" yaml:"port" env:"SERVER_PORT" envDefault:"8081"`
		Secure bool   `json:"secure" yaml:"secure" env:"SERVER_SECURE"`
	} `json:"server" yaml:"server"`
	Database struct {
		Host     string `json:"host" yaml:"host" env:"DB_HOST"`
		Port     int    `json:"port" yaml:"port" env:"DB_PORT"`
		Name     string `json:"name" yaml:"name" env:"DB_NAME"`
		User     string `json:"user" yaml:"user" env:"DB_USER"`
		Password string `json:"password" yaml:"password" env:"DB_PASSWORD"`
	} `json:"database" yaml:"database"`
}

func New(opts ...func(*ServerConfig)) *ServerConfig {
	config := &ServerConfig{}
	for _, o := range opts {
		o(config)
	}
	return config
}

func WithEnv() func(*ServerConfig) {
	return func(config *ServerConfig) {
		if err := cfg.UpdateConfigFromEnv(config); err != nil {
			log.Error(err)
		}
	}
}

func WithFile() func(config *ServerConfig) {
	return func(config *ServerConfig) {
		var fCfg ServerConfig
		if err := cfg.UpdateConfigFromFile(config, &fCfg, config.File); err != nil {
			log.Error(err)
		}
	}
}

func (c *ServerConfig) GetRepoURL() string {
	if c.Database.Host == "" {
		return ""
	}

	db := c.Database
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		db.User, db.Password, db.Host, db.Port, db.Name)
}

func (c *ServerConfig) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

func (c *ServerConfig) IsServerSecure() bool {
	return c.Server.Secure
}
