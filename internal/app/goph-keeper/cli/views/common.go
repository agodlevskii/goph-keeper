package views

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
)

type viewer interface {
	getItem() error
	getItems() error
	saveItem() error
	deleteItem() error
}

type MenuOption string

const (
	MBinary   MenuOption = "Binaries"
	MCard     MenuOption = "Cards"
	MPassword MenuOption = "Passwords"
	MText     MenuOption = "Texts"
	MExit     MenuOption = "Exit"
)

type commandOption string

const (
	cGet    commandOption = "Get a single item by ID"
	cGetAll commandOption = "Get the list of items"
	cSave   commandOption = "Store new item"
	cDelete commandOption = "Delete the existing item"
	cBack   commandOption = "Back to main menu"
)

var (
	MenuList     = []MenuOption{MBinary, MCard, MPassword, MText, MExit}
	commandList  = []commandOption{cGet, cGetAll, cSave, cDelete, cBack}
	commonHeader = []string{"ID", "Name", "Data", "Note"}
)

func getOptionsMenu(opt MenuOption) (commandOption, error) {
	mp := promptui.Select{
		Label: fmt.Sprintf("What would you like to do with %s?", strings.ToLower(string(opt))),
		Items: commandList,
	}

	_, res, err := mp.Run()
	return commandOption(res), err
}

func showMenu(v viewer, opt MenuOption) error {
	cmd, err := getOptionsMenu(opt)
	if err != nil {
		return err
	}

	switch cmd {
	case cGet:
		err = v.getItem()
	case cGetAll:
		err = v.getItems()
	case cSave:
		err = v.saveItem()
	case cDelete:
		err = v.deleteItem()
	case cBack:
		return nil
	}

	if err != nil {
		log.Error(err)
	}
	return showMenu(v, opt)
}

func getCtxTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*30)
}
