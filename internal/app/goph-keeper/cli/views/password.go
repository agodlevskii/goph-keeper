package views

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/cli/inputs"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/client"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/models"
)

type Password struct {
	keeper client.PasswordClient
}

var passwordHeader = []string{"ID", "Name", "User", "Password", "Note"}

func NewPasswordView(keeper client.PasswordClient) *Password {
	return &Password{keeper: keeper}
}

func (v *Password) ShowMenu() error {
	return showMenu(v, MPassword)
}

func (v *Password) getItem() error {
	id, err := inputs.ItemID()
	if err != nil {
		return err
	}

	data, err := v.keeper.GetPasswordByID(id)
	if err != nil {
		return err
	}

	v.showItems([]models.PasswordResponse{data})
	return nil
}

func (v *Password) getItems() error {
	items, err := v.keeper.GetAllPasswords()
	if err != nil {
		return err
	}
	v.showItems(items)
	return nil
}

func (v *Password) saveItem() error {
	name, err := inputs.ItemName()
	if err != nil {
		return err
	}

	user, err := inputs.Username()
	if err != nil {
		return err
	}

	password, err := inputs.Password()
	if err != nil {
		return err
	}

	note, err := inputs.ItemNote()
	if err != nil {
		return err
	}

	_, err = v.keeper.StorePassword(name, user, password, note)
	return err
}

func (v *Password) deleteItem() error {
	id, err := inputs.ItemID()
	if err != nil {
		return err
	}

	if err = v.keeper.DeletePassword(id); err != nil {
		return err
	}
	fmt.Print("Password item has been deleted successfully.")
	return err
}

func (v *Password) showItems(items []models.PasswordResponse) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(passwordHeader)
	for _, item := range items {
		table.Append(item.TableRow())
	}
	table.Render()
}
