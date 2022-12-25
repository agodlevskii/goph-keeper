package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/cli"
)

func main() {
	client, err := cli.NewCLI()
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Start(); err != nil {
		log.Fatal(err)
	}
}
