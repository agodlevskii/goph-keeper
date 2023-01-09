package client

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/client/config"
)

func TestNewClient(t *testing.T) {
	type args struct {
		cfg *config.ClientConfig
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "HTTP client init",
			args: args{cfg: config.New()},
			want: "client.HTTPKeeperClient",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.cfg)
			rGot := reflect.ValueOf(got)
			assert.Equal(t, tt.want, rGot.Type().String())
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
