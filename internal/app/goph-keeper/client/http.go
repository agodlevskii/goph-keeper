package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/agodlevskii/goph-keeper/internal/pkg/services/auth"

	log "github.com/sirupsen/logrus"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/client/config"
	"github.com/agodlevskii/goph-keeper/internal/pkg/cert"
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
	var (
		body []byte
		res  *http.Response
		err  error
	)

	body, err = json.Marshal(auth.Payload{
		Name:     user,
		Password: password,
	})
	if err != nil {
		return err
	}

	res, err = makeRequest(c.http, c.apiURL.String()+"/auth/login", body)
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
	var (
		res *http.Response
		err error
	)

	res, err = makeRequest(c.http, c.apiURL.String()+"/auth/logout", []byte{})
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
	var (
		body []byte
		res  *http.Response
		err  error
	)

	body, err = json.Marshal(auth.Payload{
		Name:     user,
		Password: password,
	})
	if err != nil {
		return err
	}

	res, err = makeRequest(c.http, c.apiURL.String()+"/auth/register", body)
	if err != nil {
		return err
	}
	defer closeResponseBody(res)

	if res.StatusCode != 200 {
		return errors.New("error")
	}
	return nil
}

func getClientConfig() KeeperClientConfig {
	return config.New(config.WithEnv(), config.WithFile())
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
