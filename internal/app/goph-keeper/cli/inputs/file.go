package inputs

import (
	"github.com/manifoldco/promptui"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/validators"
)

func FilePath() (string, error) {
	pp := promptui.Prompt{Label: "Enter the file path", Validate: validators.Min(5)}
	return pp.Run()
}
