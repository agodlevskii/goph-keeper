package cert

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"path/filepath"
)

type CertificatePaths struct {
	ca     CertificatePath
	client CertificatePath
	Server CertificatePath
}

type CertificatePath struct {
	Key  string
	Cert string
}

func GetCertificatePool() (*x509.CertPool, error) {
	cPaths, err := GetCertificatePaths()
	if err != nil {
		return &x509.CertPool{}, err
	}

	caCert, err := os.ReadFile(cPaths.ca.Cert)
	if err != nil {
		return &x509.CertPool{}, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	return caCertPool, nil
}

func GetClientCertificate() (tls.Certificate, error) {
	cPaths, err := GetCertificatePaths()
	if err != nil {
		return tls.Certificate{}, err
	}
	return tls.LoadX509KeyPair(cPaths.client.Cert, cPaths.client.Key)
}

func GetCertificatePaths() (CertificatePaths, error) {
	certDir, err := getCertDirPath()
	if err != nil {
		return CertificatePaths{}, err
	}
	return CertificatePaths{
		ca: CertificatePath{Cert: filepath.Join(certDir, "ca.crt")},
		client: CertificatePath{
			Key:  filepath.Join(certDir, "cli.Key"),
			Cert: filepath.Join(certDir, "cli.crt"),
		},
		Server: CertificatePath{
			Key:  filepath.Join(certDir, "Server.Key"),
			Cert: filepath.Join(certDir, "Server.crt"),
		},
	}, nil
}

func getCertDirPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(wd, "..", "..", "..", "Cert"), nil
}
