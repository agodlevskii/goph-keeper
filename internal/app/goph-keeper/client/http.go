package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/services/auth"
	"github.com/agodlevskii/goph-keeper/internal/pkg/cfg/client_config"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/agodlevskii/goph-keeper/internal/pkg/cert"
)

type HTTPKeeperClient struct {
	http   *http.Client
	apiURL string
}

func NewHTTPClient() (HTTPKeeperClient, error) {
	caCertPool, err := cert.GetCertificatePool()
	if err != nil {
		return HTTPKeeperClient{}, err
	}

	c, err := cert.GetClientCertificate()
	if err != nil {
		return HTTPKeeperClient{}, err
	}

	cfg := getClientConfig()
	return HTTPKeeperClient{
		http: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:      caCertPool,
					Certificates: []tls.Certificate{c},
				},
			},
		},
		apiURL: cfg.GetAPIAddress(),
	}, nil
}

func (c HTTPKeeperClient) Login(user, password string) {
	var (
		body []byte
		res  *http.Response
		err  error
	)

	body, err = json.Marshal(auth.Request{
		Name:     user,
		Password: password,
	})
	if err != nil {
		return
	}

	res, err = makeRequest(c.http, c.apiURL+"/auth/login", body)
	if err != nil {
		log.Error(err)
		return
	}
	defer closeResponseBody(res)

	token, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info(string(token))
}

func getClientConfig() KeeperClientConfig {
	return client_config.New(client_config.WithEnv(), client_config.WithFile())
}

func makeRequest(client *http.Client, url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	return client.Do(req)
}

func closeResponseBody(r *http.Response) {
	if err := r.Body.Close(); err != nil {
		log.Error(err)
	}
}
