package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/config"

	log "github.com/sirupsen/logrus"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/handlers"
)

var (
	buildVersion string
	buildDate    string
)

type ServerConfig interface {
	GetRepoURL() string
	GetServerAddress() string
	IsServerSecure() bool
	GetCACertPool() (*x509.CertPool, error)
	GetCertificatePaths() []string
}

func main() {
	printCompilationInfo()
	cfg := config.New(config.WithEnv(), config.WithFile())
	s, err := getServer(cfg)
	if err != nil {
		log.Fatal(err)
	}
	idleConnectionsClosed := make(chan any)

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTERM)
		<-exit
		stopServer(s)
		close(idleConnectionsClosed)
	}()

	go startServer(s, cfg)
	<-idleConnectionsClosed
}

func getServer(cfg ServerConfig) (*http.Server, error) {
	h, err := handlers.NewHandler(cfg.GetRepoURL())
	if err != nil {
		return nil, err
	}

	s := &http.Server{
		Addr:              cfg.GetServerAddress(),
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if cfg.IsServerSecure() {
		tlcfg, tErr := getTLSConfig(cfg)
		if tErr != nil {
			return nil, tErr
		}
		s.TLSConfig = tlcfg
	}

	return s, nil
}

func startServer(s *http.Server, cfg ServerConfig) {
	var err error
	if cfg.IsServerSecure() {
		paths := cfg.GetCertificatePaths()
		err = s.ListenAndServeTLS(paths[0], paths[1])
	} else {
		err = s.ListenAndServe()
	}

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func stopServer(s *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error(err)
	}
}

func getTLSConfig(cfg ServerConfig) (*tls.Config, error) {
	caCertPool, err := cfg.GetCACertPool()
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
		MinVersion: tls.VersionTLS12,
	}, nil
}

func printCompilationInfo() {
	version := getCompilationInfoValue(buildVersion)
	date := getCompilationInfoValue(buildDate)
	fmt.Printf("Build version: %s\nBuild date: %s\n\n", version, date)
}

func getCompilationInfoValue(v string) string {
	if v != "" {
		return v
	}
	return "N/A"
}
