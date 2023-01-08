package inputs

import (
	"github.com/manifoldco/promptui"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/validators"
)

func Username() (string, error) {
	up := promptui.Prompt{Label: "Enter the username", Validate: validators.Min(1)}
	return up.Run()
}

func Password() (string, error) {
	pp := promptui.Prompt{Label: "Enter the user password", Validate: validators.Min(1)}
	return pp.Run()
}

func LoginRetry() (string, error) {
	ep := promptui.Prompt{
		Label:    "Incorrect name or password. Would you like to try again? (y/N)",
		Validate: validators.Min(1),
	}
	return ep.Run()
}
