package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/handlers"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/storage"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	buildVersion string
	buildDate    string
)

func main() {
	printCompilationInfo()
	s := getServer()
	idleConnectionsClosed := make(chan any)

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTERM)
		<-exit
		stopServer(s)
		close(idleConnectionsClosed)
	}()

	go startServer(s)
	<-idleConnectionsClosed
}

func getServer() *http.Server {
	db, err := storage.NewStorage("")
	if err != nil {
		log.Fatal(err)
	}

	h := handlers.NewHandler(db)
	return &http.Server{
		Addr:              "localhost:8081",
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func startServer(s *http.Server) {
	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error(err)
	}
}

func stopServer(s *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error(err)
	}
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
