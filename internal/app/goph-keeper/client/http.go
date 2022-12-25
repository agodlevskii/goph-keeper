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
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/models"
	"github.com/agodlevskii/goph-keeper/internal/pkg/cert"
)

type HTTPKeeperClient struct {
	http   *http.Client
	apiURL *url.URL
}

const (
	SBinary   string = "/storage/binary/"
	SCard     string = "/storage/card/"
	SPassword string = "/storage/password/"
	SText     string = "/storage/text/"
)

var ErrUnauthorized = errors.New("incorrect username or password")

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
	res, err := c.makeRequest(http.MethodPost, "/auth/login", models.UserRequest{
		Name:     user,
		Password: password,
	})
	if err != nil {
		if res.StatusCode == http.StatusUnauthorized {
			return ErrUnauthorized
		}
		return err
	}
	defer closeResponseBody(res.Body)

	c.http.Jar.SetCookies(c.apiURL, res.Cookies())
	return nil
}

func (c HTTPKeeperClient) Logout() error {
	res, err := c.makeRequest(http.MethodPost, "/auth/logout", nil)
	if err != nil {
		return err
	}
	defer closeResponseBody(res.Body)

	c.http.Jar.SetCookies(c.apiURL, nil)
	return nil
}

func (c HTTPKeeperClient) Register(user, password string) error {
	res, err := c.makeRequest(http.MethodPost, "/auth/register", models.UserRequest{
		Name:     user,
		Password: password,
	})
	if err != nil {
		return err
	}
	defer closeResponseBody(res.Body)
	return nil
}

func (c HTTPKeeperClient) DeleteBinary(id string) error {
	return c.deleteData(SBinary, id)
}

func (c HTTPKeeperClient) GetAllBinaries() ([]models.BinaryResponse, error) {
	body, err := c.getAllData(SBinary)
	if err != nil {
		return nil, err
	}
	defer closeResponseBody(body)

	var bins []models.BinaryResponse
	err = json.NewDecoder(body).Decode(&bins)
	return bins, err
}

func (c HTTPKeeperClient) GetBinaryByID(id string) (models.BinaryResponse, error) {
	var data models.BinaryResponse
	body, err := c.getDataByID(SBinary, id)
	if err != nil {
		return data, err
	}
	defer closeResponseBody(body)

	err = json.NewDecoder(body).Decode(&data)
	return data, err
}

func (c HTTPKeeperClient) StoreBinary(name string, data []byte, note string) (string, error) {
	return c.storeData(SBinary, models.BinaryRequest{
		Name: name,
		Data: data,
		Note: note,
	})
}

func (c HTTPKeeperClient) DeleteText(id string) error {
	return c.deleteData(SText, id)
}

func (c HTTPKeeperClient) GetAllTexts() ([]models.TextResponse, error) {
	body, err := c.getAllData(SText)
	if err != nil {
		return nil, err
	}
	defer closeResponseBody(body)

	var data []models.TextResponse
	err = json.NewDecoder(body).Decode(&data)
	return data, err
}

func (c HTTPKeeperClient) GetTextByID(id string) (models.TextResponse, error) {
	var data models.TextResponse
	body, err := c.getDataByID(SText, id)
	if err != nil {
		return data, err
	}
	defer closeResponseBody(body)

	err = json.NewDecoder(body).Decode(&data)
	return data, err
}

func (c HTTPKeeperClient) StoreText(name, data, note string) (string, error) {
	return c.storeData(SText, models.TextRequest{
		Name: name,
		Data: data,
		Note: note,
	})
}

func (c HTTPKeeperClient) DeleteCard(id string) error {
	return c.deleteData(SCard, id)
}

func (c HTTPKeeperClient) GetAllCards() ([]models.CardResponse, error) {
	body, err := c.getAllData(SCard)
	if err != nil {
		return nil, err
	}
	defer closeResponseBody(body)

	var data []models.CardResponse
	err = json.NewDecoder(body).Decode(&data)
	return data, err
}

func (c HTTPKeeperClient) GetCardByID(id string) (models.CardResponse, error) {
	var data models.CardResponse
	body, err := c.getDataByID(SCard, id)
	if err != nil {
		return data, err
	}
	defer closeResponseBody(body)

	err = json.NewDecoder(body).Decode(&data)
	return data, err
}

func (c HTTPKeeperClient) StoreCard(name, number, holder, expDate, cvv, note string) (string, error) {
	return c.storeData(SCard, models.CardRequest{
		Name:    name,
		Number:  number,
		Holder:  holder,
		ExpDate: expDate,
		CVV:     cvv,
		Note:    note,
	})
}

func (c HTTPKeeperClient) DeletePassword(id string) error {
	return c.deleteData(SPassword, id)
}

func (c HTTPKeeperClient) GetAllPasswords() ([]models.PasswordResponse, error) {
	body, err := c.getAllData(SPassword)
	if err != nil {
		return nil, err
	}
	defer closeResponseBody(body)

	var data []models.PasswordResponse
	err = json.NewDecoder(body).Decode(&data)
	return data, err
}

func (c HTTPKeeperClient) GetPasswordByID(id string) (models.PasswordResponse, error) {
	var data models.PasswordResponse
	body, err := c.getDataByID(SPassword, id)
	if err != nil {
		return data, err
	}
	defer closeResponseBody(body)

	err = json.NewDecoder(body).Decode(&data)
	return data, err
}

func (c HTTPKeeperClient) StorePassword(name, user, password, note string) (string, error) {
	return c.storeData(SPassword, models.PasswordRequest{
		Name:     name,
		User:     user,
		Password: password,
		Note:     note,
	})
}

func (c HTTPKeeperClient) deleteData(url, id string) error {
	_, err := c.makeRequest(http.MethodDelete, url+id, nil)
	return err
}

func (c HTTPKeeperClient) getAllData(url string) (io.ReadCloser, error) {
	res, err := c.makeRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

func (c HTTPKeeperClient) getDataByID(url, id string) (io.ReadCloser, error) {
	res, err := c.makeRequest(http.MethodGet, url+id, nil)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

func (c HTTPKeeperClient) storeData(url string, data any) (string, error) {
	res, err := c.makeRequest(http.MethodPost, url, data)
	if err != nil {
		return "", err
	}
	defer closeResponseBody(res.Body)

	id, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(id), err
}

func (c HTTPKeeperClient) makeRequest(method, url string, data any) (*http.Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var reader io.Reader
	if len(body) > 0 {
		reader = bytes.NewBuffer(body)
	}

	req, err := http.NewRequest(method, c.apiURL.String()+url, reader)
	if err != nil {
		return nil, err
	}

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return res, errors.New("response code")
	}
	return res, err
}

func getClientConfig() KeeperClientConfig {
	return config.New(config.WithEnv(), config.WithFile())
}

func closeResponseBody(b io.Closer) {
	if err := b.Close(); err != nil {
		log.Error(err)
	}
}
