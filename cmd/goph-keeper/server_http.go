package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/handlers"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/storage"
	"github.com/agodlevskii/goph-keeper/internal/pkg/cert"
	"github.com/agodlevskii/goph-keeper/internal/pkg/cfg/server_config"
)

var (
	buildVersion string
	buildDate    string
)

type ServerConfig interface {
	GetRepoURL() string
	GetServerAddress() string
	IsServerSecure() bool
}

func main() {
	printCompilationInfo()
	sCfg := server_config.New(server_config.WithEnv(), server_config.WithFile())
	s, err := getServer(sCfg)
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

	go startServer(s, sCfg)
	<-idleConnectionsClosed
}

func getServer(sCfg ServerConfig) (*http.Server, error) {
	db, err := storage.NewStorage(sCfg.GetRepoURL())
	if err != nil {
		return nil, err
	}

	h := handlers.NewHandler(db)
	s := &http.Server{
		Addr:              sCfg.GetServerAddress(),
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if sCfg.IsServerSecure() {
		tlsCfg, tErr := getTLSConfig()
		if tErr != nil {
			return nil, tErr
		}
		s.TLSConfig = tlsCfg
	}

	return s, nil
}

func startServer(s *http.Server, sCfg ServerConfig) {
	cPaths, err := cert.GetCertificatePaths()
	if err != nil {
		log.Fatal(err)
	}

	if sCfg.IsServerSecure() {
		err = s.ListenAndServeTLS(cPaths.Server.Cert, cPaths.Server.Key)
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

func getTLSConfig() (*tls.Config, error) {
	caCertPool, err := cert.GetCertificatePool()
	if err != nil {
		return &tls.Config{}, err
	}

	return &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
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
