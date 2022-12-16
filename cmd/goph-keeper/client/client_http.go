package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/client"
)

func main() {
	c, err := client.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	c.Login("test", "super")
}
