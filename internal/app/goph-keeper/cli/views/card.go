package views

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/cli/inputs"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/client"
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/models"
)

type Card struct {
	keeper client.CardClient
}

var cardHeader = []string{"ID", "Name", "Number", "Holder", "Expire date", "CardCVV", "Note"}

func NewCardView(keeper client.CardClient) *Card {
	return &Card{keeper: keeper}
}

func (v *Card) ShowMenu() error {
	return showMenu(v, MCard)
}

func (v *Card) getItem() error {
	id, err := inputs.ItemID()
	if err != nil {
		return err
	}

	ctx, cancel := getCtxTimeout()
	defer cancel()

	data, err := v.keeper.GetCardByID(ctx, id)
	if err != nil {
		return err
	}

	v.showItems([]models.CardResponse{data})
	return nil
}

func (v *Card) getItems() error {
	ctx, cancel := getCtxTimeout()
	defer cancel()

	items, err := v.keeper.GetAllCards(ctx)
	if err != nil {
		return err
	}
	v.showItems(items)
	return nil
}

func (v *Card) saveItem() error {
	name, err := inputs.ItemName()
	if err != nil {
		return err
	}
	number, err := inputs.CardNumber()
	if err != nil {
		return err
	}
	holder, err := inputs.CardHolder()
	if err != nil {
		return err
	}
	expDate, err := inputs.CardExpDate()
	if err != nil {
		return err
	}
	cvv, err := inputs.CardCVV()
	if err != nil {
		return err
	}
	note, err := inputs.ItemNote()
	if err != nil {
		return err
	}

	ctx, cancel := getCtxTimeout()
	defer cancel()

	_, err = v.keeper.StoreCard(ctx, name, number, holder, expDate, cvv, note)
	return err
}

func (v *Card) deleteItem() error {
	id, err := inputs.ItemID()
	if err != nil {
		return err
	}

	ctx, cancel := getCtxTimeout()
	defer cancel()

	if err = v.keeper.DeleteCard(ctx, id); err != nil {
		return err
	}
	fmt.Print("Card item has been deleted successfully.")
	return err
}

func (v *Card) showItems(items []models.CardResponse) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(cardHeader)
	for _, item := range items {
		table.Append(item.TableRow())
	}
	table.Render()
}
