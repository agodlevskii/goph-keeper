package client

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"net/http"
	"strconv"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/client/config"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/models"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/handlers"
)

var (
	s http.Server
	c HTTPKeeperClient
)

func TestHTTPKeeperClient_DeleteBinary(t *testing.T) {
	type args struct {
		id   string
		item models.BinaryRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "No ID passed",
			wantErr: errors.New("405 Method Not Allowed"),
		},
		{
			name:    "ID doesn't exist",
			args:    args{id: "test"},
			wantErr: errors.New("404 Not Found"),
		},
		{
			name: "ID exists",
			args: args{id: "test", item: models.BinaryRequest{Name: "test", Data: []byte("test")}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.item.Name != "" {
				id, err := c.StoreBinary(context.Background(), tt.args.item.Name, tt.args.item.Data, "")
				if err != nil {
					t.Fatal(err)
				}
				tt.args.id = id
			}

			assert.Equal(t, tt.wantErr, c.DeleteBinary(context.Background(), tt.args.id))
		})
	}
}

func TestHTTPKeeperClient_GetAllBinaries(t *testing.T) {
	tests := []struct {
		name    string
		data    models.BinaryRequest
		want    []models.BinaryResponse
		wantErr error
	}{
		{
			name: "No data found",
			want: []models.BinaryResponse{},
		},
		{
			name: "Data found",
			data: models.BinaryRequest{Name: "test", Data: []byte("test")},
			want: []models.BinaryResponse{{Name: "test", Data: []byte{}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.data.Name != "" {
				id, err := c.StoreBinary(context.Background(), tt.data.Name, tt.data.Data, "")
				if err != nil {
					t.Fatal(err)
				}
				tt.want[0].ID = id
			}

			got, err := c.GetAllBinaries(context.Background())
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestHTTPKeeperClient_GetAllCards(t *testing.T) {
	tests := []struct {
		name    string
		data    models.CardRequest
		want    []models.CardResponse
		wantErr error
	}{
		{
			name: "No data found",
			want: []models.CardResponse{},
		},
		{
			name: "Data found",
			data: models.CardRequest{Name: "test", Number: "123"},
			want: []models.CardResponse{{Name: "test", Number: "123", CVV: "***"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.data.Name != "" {
				id, err := c.StoreCard(context.Background(), tt.data.Name, tt.data.Number, "", "", "", "")
				if err != nil {
					t.Fatal(err)
				}
				tt.want[0].ID = id
			}

			got, err := c.GetAllCards(context.Background())
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestHTTPKeeperClient_GetCardByID(t *testing.T) {
	type args struct {
		id   string
		item models.CardRequest
	}
	tests := []struct {
		name    string
		args    args
		want    models.CardResponse
		wantErr error
	}{
		{
			name:    "No ID found",
			args:    args{id: "test"},
			wantErr: errors.New("404 Not Found"),
		},
		{
			name: "Data found",
			args: args{id: "test", item: models.CardRequest{Name: "test", Number: "123"}},
			want: models.CardResponse{Name: "test", Number: "123"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.item.Name != "" {
				id, err := c.StoreCard(context.Background(), tt.args.item.Name, tt.args.item.Number, "", "", "", "")
				if err != nil {
					t.Fatal(err)
				}
				tt.args.id = id
				tt.want.ID = id
			}

			got, err := c.GetCardByID(context.Background(), tt.args.id)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestHTTPKeeperClient_GetBinaryByID(t *testing.T) {
	type args struct {
		id   string
		item models.BinaryRequest
	}
	tests := []struct {
		name    string
		args    args
		want    models.BinaryResponse
		wantErr error
	}{
		{
			name:    "No ID found",
			args:    args{id: "test"},
			wantErr: errors.New("404 Not Found"),
		},
		{
			name: "Data found",
			args: args{id: "test", item: models.BinaryRequest{Name: "test", Data: []byte("test")}},
			want: models.BinaryResponse{Name: "test", Data: []byte("test")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.item.Name != "" {
				id, err := c.StoreBinary(context.Background(), tt.args.item.Name, tt.args.item.Data, "")
				if err != nil {
					t.Fatal(err)
				}
				tt.args.id = id
				tt.want.ID = id
			}

			got, err := c.GetBinaryByID(context.Background(), tt.args.id)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestHTTPKeeperClient_StoreBinary(t *testing.T) {
	type args struct {
		name string
		data []byte
		note string
	}
	tests := []struct {
		name      string
		args      args
		wantEmpty bool
		wantErr   error
	}{
		{
			name:      "Empty payload",
			wantEmpty: true,
			wantErr:   errors.New("400 Bad Request"),
		},
		{
			name: "Correct payload",
			args: args{
				name: "test",
				data: []byte("test"),
				note: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.StoreBinary(context.Background(), tt.args.name, tt.args.data, tt.args.note)
			assert.Equal(t, tt.wantEmpty, got == "")
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestHTTPKeeperClient_GetAllPasswords(t *testing.T) {
	tests := []struct {
		name    string
		data    models.PasswordRequest
		want    []models.PasswordResponse
		wantErr error
	}{
		{
			name: "No data found",
			want: []models.PasswordResponse{},
		},
		{
			name: "Data found",
			data: models.PasswordRequest{Name: "test", User: "test", Password: "test"},
			want: []models.PasswordResponse{{Name: "test", User: "test", Password: "********"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.data.Name != "" {
				id, err := c.StorePassword(context.Background(), tt.data.Name, tt.data.User, tt.data.Password, "")
				if err != nil {
					t.Fatal(err)
				}
				tt.want[0].ID = id
			}

			got, err := c.GetAllPasswords(context.Background())
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestHTTPKeeperClient_GetAllTexts(t *testing.T) {
	tests := []struct {
		name    string
		data    models.TextRequest
		want    []models.TextResponse
		wantErr error
	}{
		{
			name: "No data found",
			want: []models.TextResponse{},
		},
		{
			name: "Data found",
			data: models.TextRequest{Name: "test", Data: "test"},
			want: []models.TextResponse{{Name: "test"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.data.Name != "" {
				id, err := c.StoreText(context.Background(), tt.data.Name, tt.data.Data, "")
				if err != nil {
					t.Fatal(err)
				}
				tt.want[0].ID = id
			}

			got, err := c.GetAllTexts(context.Background())
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestHTTPKeeperClient_GetTextByID(t *testing.T) {
	type args struct {
		id   string
		item models.TextRequest
	}
	tests := []struct {
		name    string
		args    args
		want    models.TextResponse
		wantErr error
	}{
		{
			name:    "No ID found",
			args:    args{id: "test"},
			wantErr: errors.New("404 Not Found"),
		},
		{
			name: "Data found",
			args: args{id: "test", item: models.TextRequest{Name: "test", Data: "test"}},
			want: models.TextResponse{Name: "test", Data: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.item.Name != "" {
				id, err := c.StoreText(context.Background(), tt.args.item.Name, tt.args.item.Data, "")
				if err != nil {
					t.Fatal(err)
				}
				tt.args.id = id
				tt.want.ID = id
			}

			got, err := c.GetTextByID(context.Background(), tt.args.id)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestHTTPKeeperClient_Login(t *testing.T) {
	type args struct {
		user     string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "No credential passed",
			wantErr: errors.New("400 Bad Request"),
		},
		{
			name: "Correct credential passed",
			args: args{user: "test1", password: "test1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.Login(context.Background(), tt.args.user, tt.args.password)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestHTTPKeeperClient_Logout(t *testing.T) {
	tests := []struct {
		name    string
		wantErr error
	}{
		{
			name: "Logout request",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.Logout(context.Background())
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestHTTPKeeperClient_Register(t *testing.T) {
	type args struct {
		user     string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "Empty payload",
			wantErr: errors.New("400 Bad Request"),
		},
		{
			name: "Correct payload",
			args: args{user: "test", password: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.Register(context.Background(), tt.args.user, tt.args.password)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func init() {
	port, err := generateRandomPort()
	if err != nil {
		log.Fatal(err)
	}
	go startServer(port)

	c, err = initTestClient(port)
	if err != nil {
		log.Fatal(err)
	}

	if err = c.Register(context.Background(), "user", "user"); err != nil {
		log.Fatal(err)
	}
	if err = c.Login(context.Background(), "user", "user"); err != nil {
		log.Fatal(err)
	}
}

func startServer(port int) {
	h, err := handlers.NewHandler("")
	if err != nil {
		log.Fatal(err)
	}

	s = http.Server{
		Addr:              ":" + strconv.Itoa(port),
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
	}
	if err = s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func initTestClient(port int) (HTTPKeeperClient, error) {
	cfg := config.New()
	cfg.API.Host = "localhost"
	cfg.API.Port = port
	cfg.API.Route = "/api/v1"
	return NewHTTPClient(cfg)
}

func generateRandomPort() (int, error) {
	r, err := rand.Int(rand.Reader, big.NewInt(9999))
	if err != nil {
		return 0, err
	}
	return int(r.Int64()), nil
}
