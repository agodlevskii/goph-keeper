package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/client/config"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/models"
	"github.com/agodlevskii/goph-keeper/internal/pkg/cert"
	"github.com/agodlevskii/goph-keeper/internal/pkg/services/auth"
)

type HTTPKeeperClient struct {
	http   *http.Client
	apiURL *url.URL
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

	jar, err := cookiejar.New(nil)
	if err != nil {
		return HTTPKeeperClient{}, err
	}

	cfg := getClientConfig()
	uri, err := url.Parse(cfg.GetAPIAddress())
	if err != nil {
		return HTTPKeeperClient{}, err
	}

	return HTTPKeeperClient{
		http: &http.Client{
			Jar: jar,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:      caCertPool,
					Certificates: []tls.Certificate{c},
					MinVersion:   tls.VersionTLS12,
				},
			},
		},
		apiURL: uri,
	}, nil
}

func (c HTTPKeeperClient) Login(user, password string) error {
	body, err := json.Marshal(auth.Payload{
		Name:     user,
		Password: password,
	})
	if err != nil {
		return err
	}

	res, err := makeRequest(c.http, http.MethodPost, c.apiURL.String()+"/auth/login", body)
	if err != nil {
		return err
	}
	defer closeResponseBody(res)

	if res.StatusCode != 200 {
		return errors.New("error")
	}

	c.http.Jar.SetCookies(c.apiURL, res.Cookies())
	return nil
}

func (c HTTPKeeperClient) Logout() error {
	res, err := makeRequest(c.http, http.MethodPost, c.apiURL.String()+"/auth/logout", []byte{})
	if err != nil {
		return err
	}
	defer closeResponseBody(res)

	if res.StatusCode != 200 {
		return errors.New("error")
	}

	c.http.Jar.SetCookies(c.apiURL, nil)
	return nil
}

func (c HTTPKeeperClient) Register(user, password string) error {
	body, err := json.Marshal(models.UserRequest{
		Name:     user,
		Password: password,
	})
	if err != nil {
		return err
	}

	res, err := makeRequest(c.http, http.MethodPost, c.apiURL.String()+"/auth/register", body)
	if err != nil {
		return err
	}
	defer closeResponseBody(res)

	if res.StatusCode != 200 {
		return errors.New("error")
	}
	return nil
}

func (c HTTPKeeperClient) DeleteBinary(id string) error {
	res, err := makeRequest(c.http, http.MethodDelete, c.apiURL.String()+"/storage/binary/"+id, nil)
	if err != nil {
		return err
	}
	defer closeResponseBody(res)

	if res.StatusCode != 200 {
		return errors.New("error")
	}
	return nil
}

func (c HTTPKeeperClient) GetAllBinaries() ([]models.BinaryResponse, error) {
	res, err := makeRequest(c.http, http.MethodGet, c.apiURL.String()+"/storage/binary/", nil)
	if err != nil {
		return nil, err
	}
	defer closeResponseBody(res)

	if res.StatusCode != 200 {
		return nil, errors.New("error")
	}

	var bins []models.BinaryResponse
	err = json.NewDecoder(res.Body).Decode(&bins)
	return bins, err
}

func (c HTTPKeeperClient) GetBinaryByID(id string) (models.BinaryResponse, error) {
	res, err := makeRequest(c.http, http.MethodGet, c.apiURL.String()+"/storage/binary/"+id, nil)
	if err != nil {
		return models.BinaryResponse{}, err
	}
	defer closeResponseBody(res)

	if res.StatusCode != 200 {
		return models.BinaryResponse{}, errors.New("error")
	}

	var bin models.BinaryResponse
	err = json.NewDecoder(res.Body).Decode(&bin)
	return bin, err
}

func (c HTTPKeeperClient) StoreBinary(name string, data []byte, note string) (string, error) {
	body, err := json.Marshal(models.BinaryRequest{
		Name: name,
		Data: data,
		Note: note,
	})
	if err != nil {
		return "", err
	}

	res, err := makeRequest(c.http, http.MethodPost, c.apiURL.String()+"/storage/binary/", body)
	if err != nil {
		return "", err
	}
	defer closeResponseBody(res)

	if res.StatusCode != 200 {
		return "", errors.New("error")
	}

	result, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(result), err
}

func getClientConfig() KeeperClientConfig {
	return config.New(config.WithEnv(), config.WithFile())
}

func makeRequest(client *http.Client, method string, url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
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
