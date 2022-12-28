package cli

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/cli/inputs"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/cli/views"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/client"
)

type View interface {
	ShowMenu() error
}

type AppCLI struct {
	client   client.KeeperClient
	binary   View
	card     View
	password View
	text     View
}

func NewCLI() (*AppCLI, error) {
	c, err := client.NewClient()
	if err != nil {
		return nil, err
	}

	return &AppCLI{
		client:   c,
		binary:   views.NewBinaryView(c),
		card:     views.NewCardView(c),
		password: views.NewPasswordView(c),
		text:     views.NewTextView(c),
	}, nil
}

func (app *AppCLI) Start() error {
	if err := app.login(); err != nil {
		if errors.Is(err, client.ErrUnauthorized) {
			if retry, rErr := inputs.LoginRetry(); rErr != nil {
				return rErr
			} else if strings.ToLower(retry)[:1] != "n" {
				return app.Start()
			}
		}
		return err
	}
	return app.mainMenu()
}

func (app *AppCLI) login() error {
	user, err := inputs.Username()
	if err != nil {
		return err
	}

	password, err := inputs.Password()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	return app.client.Login(ctx, user, password)
}

func (app *AppCLI) mainMenu() error {
	mp := promptui.Select{
		Label: "What type of data would you like to work with?",
		Items: views.MenuList,
	}

	_, opt, err := mp.Run()
	if err != nil {
		return err
	}

	switch views.MenuOption(opt) {
	case views.MBinary:
		err = app.binary.ShowMenu()
	case views.MCard:
		err = app.card.ShowMenu()
	case views.MPassword:
		err = app.password.ShowMenu()
	case views.MText:
		err = app.text.ShowMenu()
	case views.MExit:
		return nil
	}

	if err != nil {
		log.Error(err)
	}
	return app.mainMenu()
}
