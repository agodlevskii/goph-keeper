package views

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/olekukonko/tablewriter"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/cli/inputs"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/client"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/models"
)

type Binary struct {
	keeper client.BinaryClient
}

func NewBinaryView(keeper client.BinaryClient) *Binary {
	return &Binary{keeper: keeper}
}

func (v *Binary) ShowMenu() error {
	return showMenu(v, MBinary)
}

func (v *Binary) getItem() error {
	id, err := inputs.ItemID()
	if err != nil {
		return err
	}

	data, err := v.keeper.GetBinaryByID(id)
	if err != nil {
		return err
	}

	v.showItems([]models.BinaryResponse{data})
	return nil
}

func (v *Binary) getItems() error {
	items, err := v.keeper.GetAllBinaries()
	if err != nil {
		return err
	}
	v.showItems(items)
	return nil
}

func (v *Binary) saveItem() error {
	path, err := inputs.FilePath()
	if err != nil {
		return err
	}
	_, name := filepath.Split(path)

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	note, err := inputs.ItemNote()
	if err != nil {
		return err
	}

	_, err = v.keeper.StoreBinary(name, data, note)
	return err
}

func (v *Binary) deleteItem() error {
	id, err := inputs.ItemID()
	if err != nil {
		return err
	}

	if err = v.keeper.DeleteBinary(id); err != nil {
		return err
	}
	fmt.Print("Binary item has been deleted successfully.")
	return err
}

func (v *Binary) showItems(items []models.BinaryResponse) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(commonHeader)
	for _, item := range items {
		table.Append(item.TableRow())
	}
	table.Render()
}
