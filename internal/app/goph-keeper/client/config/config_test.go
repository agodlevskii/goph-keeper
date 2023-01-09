package config

import (
	"crypto/tls"
	"crypto/x509"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientConfig_GetAPIAddress(t *testing.T) {
	tests := []struct {
		name string
		cfg  ClientConfig
		want string
	}{
		{
			name: "Empty config",
		},
		{
			name: "Secure connection",
			cfg: ClientConfig{
				API: struct {
					Host   string `json:"host" yaml:"host" env:"API_HOST"`
					Port   int    `json:"port" yaml:"port" env:"API_PORT" envDefault:"8081"`
					Route  string `json:"route" yaml:"route" env:"API_ROUTE" envDefault:"/api/v1"`
					Secure bool   `json:"secure" yaml:"secure" env:"CLIENT_SECURE"`
				}{Port: 8081, Secure: true},
			},
			want: "https://:8081",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.cfg.GetAPIAddress())
		})
	}
}

func TestClientConfig_GetCACertPool(t *testing.T) {
	tests := []struct {
		name    string
		cfg     ClientConfig
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

func TestClientConfig_GetCertificate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     ClientConfig
		want    tls.Certificate
		wantErr bool
	}{
		{
			name:    "Empty config",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cfg.GetCertificate()
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		opts []func(*ClientConfig)
	}
	tests := []struct {
		name string
		args args
		want *ClientConfig
	}{
		{
			name: "With env variables",
			args: args{opts: []func(config *ClientConfig){WithEnv()}},
			want: &ClientConfig{File: "client.yml"},
		},
		{
			name: "With file",
			args: args{opts: []func(config *ClientConfig){WithFile()}},
			want: &ClientConfig{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want.File, New(tt.args.opts...).File)
		})
	}
}
