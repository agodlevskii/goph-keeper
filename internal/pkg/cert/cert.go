package cert

import (
	"crypto/tls"
	"crypto/x509"
	"os"
)

func GetCertificatePool(path string) (*x509.CertPool, error) {
	caCert, err := os.ReadFile(path)
	if err != nil {
		return &x509.CertPool{}, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	return caCertPool, nil
}

func GetClientCertificate(cert string, key string) (tls.Certificate, error) {
	return tls.LoadX509KeyPair(cert, key)
}
