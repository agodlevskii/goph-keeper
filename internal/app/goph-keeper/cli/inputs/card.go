package inputs

import (
	"github.com/manifoldco/promptui"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/validators"
)

func CardNumber() (string, error) {
	cp := promptui.Prompt{Label: "Enter the card number", Validate: validators.CardNumber}
	return cp.Run()
}

func CardHolder() (string, error) {
	hp := promptui.Prompt{Label: "Enter the card holder name", Validate: validators.Max(50)}
	return hp.Run()
}

func CardExpDate() (string, error) {
	ep := promptui.Prompt{Label: "Enter the card expire date", Validate: validators.CardExpDate}
	return ep.Run()
}

func CardCVV() (string, error) {
	cp := promptui.Prompt{Label: "Enter the card's CardCVV", Validate: validators.CardCVV}
	return cp.Run()
}
