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

	id, err := c.StoreCard("test", "test", "test", "01/21", "123", "note")
	if err != nil {
		log.Error(err)
	}
	log.Info(id)

	bin, err := c.GetCardByID(id)
	if err != nil {
		log.Error(err)
	}
	log.Info(bin)

	bins, err := c.GetAllCards()
	if err != nil {
		log.Error(err)
	}
	log.Info(bins)

	if err = c.DeleteCard(id); err != nil {
		log.Error(err)
	}

	if err = c.Logout(); err != nil {
		log.Error(err)
	}
}
