package views

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/cli/inputs"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/client"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/models"
)

type Text struct {
	keeper client.TextClient
}

func NewTextView(keeper client.TextClient) *Text {
	return &Text{keeper: keeper}
}

func (v *Text) ShowMenu() error {
	return showMenu(v, MText)
}

func (v *Text) getItem() error {
	id, err := inputs.ItemID()
	if err != nil {
		return err
	}

	ctx, cancel := getCtxTimeout()
	defer cancel()

	data, err := v.keeper.GetTextByID(ctx, id)
	if err != nil {
		return err
	}

	v.showItems([]models.TextResponse{data})
	return nil
}

func (v *Text) getItems() error {
	ctx, cancel := getCtxTimeout()
	defer cancel()

	items, err := v.keeper.GetAllTexts(ctx)
	if err != nil {
		return err
	}
	v.showItems(items)
	return nil
}

func (v *Text) saveItem() error {
	name, err := inputs.ItemName()
	if err != nil {
		return err
	}
	text, err := inputs.ItemText()
	if err != nil {
		return err
	}
	note, err := inputs.ItemNote()
	if err != nil {
		return err
	}

	ctx, cancel := getCtxTimeout()
	defer cancel()

	_, err = v.keeper.StoreText(ctx, name, text, note)
	return err
}

func (v *Text) deleteItem() error {
	id, err := inputs.ItemID()
	if err != nil {
		return err
	}

	ctx, cancel := getCtxTimeout()
	defer cancel()

	if err = v.keeper.DeleteText(ctx, id); err != nil {
		return err
	}
	fmt.Print("Text item has been deleted successfully.")
	return err
}

func (v *Text) showItems(items []models.TextResponse) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(commonHeader)
	for _, item := range items {
		table.Append(item.TableRow())
	}
	table.Render()
}
