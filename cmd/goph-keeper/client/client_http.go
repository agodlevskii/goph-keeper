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

	if err = c.Login("test", "super"); err != nil {
		log.Error(err)
	}

	id, err := c.StoreBinary("test", []byte("test"), "note")
	if err != nil {
		log.Error(err)
	}
	log.Info(id)

	bin, err := c.GetBinaryByID(id)
	if err != nil {
		log.Error(err)
	}
	log.Info(bin)

	bins, err := c.GetAllBinaries()
	if err != nil {
		log.Error(err)
	}
	log.Info(bins[0].Data)

	if err = c.DeleteBinary(id); err != nil {
		log.Error(err)
	}

	if err = c.Logout(); err != nil {
		log.Error(err)
	}
}
