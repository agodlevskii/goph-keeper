package config

import (
	"crypto/x509"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []func(*ServerConfig)
	}
	tests := []struct {
		name string
		args args
		want *ServerConfig
	}{
		{
			name: "With env variables",
			args: args{opts: []func(config *ServerConfig){WithEnv()}},
			want: &ServerConfig{File: "server.yml"},
		},
		{
			name: "With file",
			args: args{opts: []func(config *ServerConfig){WithFile()}},
			want: &ServerConfig{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want.File, New(tt.args.opts...).File)
		})
	}
}

func TestServerConfig_GetCACertPool(t *testing.T) {
	tests := []struct {
		name    string
		cfg     ServerConfig
		want    *x509.CertPool
		wantErr bool
	}{
		{
			name:    "Empty config",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cfg.GetCACertPool()
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestServerConfig_GetCertificatePaths(t *testing.T) {
	tests := []struct {
		name string
		cfg  ServerConfig
		want []string
	}{
		{
			name: "Empty config",
			want: []string{"", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.cfg.GetCertificatePaths())
		})
	}
}

func TestServerConfig_GetRepoURL(t *testing.T) {
	tests := []struct {
		name string
		cfg  ServerConfig
		want string
	}{
		{
			name: "Empty config",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.cfg.GetRepoURL())
		})
	}
}

func TestServerConfig_GetServerAddress(t *testing.T) {
	tests := []struct {
		name string
		cfg  ServerConfig
		want string
	}{
		{
			name: "Empty config",
			want: ":0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.cfg.GetServerAddress())
		})
	}
}

func TestServerConfig_IsServerSecure(t *testing.T) {
	tests := []struct {
		name string
		cfg  ServerConfig
		want bool
	}{
		{
			name: "Empty config",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.cfg.IsServerSecure())
		})
	}
}

func TestWithEnv(t *testing.T) {
	tests := []struct {
		name string
		cfg  *ServerConfig
		want *ServerConfig
	}{
		{
			name: "Default config",
			cfg:  &ServerConfig{},
			want: &ServerConfig{File: "server.yml"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WithEnv()(tt.cfg)
			assert.Equal(t, tt.want.File, tt.cfg.File)
		})
	}
}
