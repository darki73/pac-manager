package prompt

import (
	"fmt"
	"github.com/darki73/pac-manager/pkg/logger"
	"github.com/manifoldco/promptui"
)

// Text prompts the user for a string.
func Text(label string) string {
	prompt := promptui.Prompt{
		Label:  label,
		Stdout: noBellStdoutInstance,
	}

	result, err := prompt.Run()

	if err != nil {
		logger.Fatalf("helpers:prompt", "failed to read input: %s", err.Error())
	}

	return result
}

// Password prompts the user for a password.
func Password(label string) string {
	prompt := promptui.Prompt{
		Label:  label,
		Stdout: noBellStdoutInstance,
		Mask:   '*',
	}

	result, err := prompt.Run()

	if err != nil {
		logger.Fatalf("helpers:prompt", "failed to read input: %s", err.Error())
	}

	return result
}

// Select is a helper function to select an item from a list.
func Select(label, selected string, items []string) string {
	prompt := promptui.Select{
		Label:  label,
		Items:  items,
		Stdout: noBellStdoutInstance,
		Templates: &promptui.SelectTemplates{
			Active:   promptui.IconSelect + " {{ . | cyan }}",
			Selected: promptui.IconGood + fmt.Sprintf(" %s {{ . | green }}", selected),
		},
	}

	_, result, err := prompt.Run()

	if err != nil {
		logger.Fatalf("helpers:prompt", "failed to read input: %s", err.Error())
	}

	return result
}

// Confirm is a helper function to confirm an action.
func Confirm(label string) bool {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
		Stdout:    noBellStdoutInstance,
		Default:   "y",
	}

	result, err := prompt.Run()

	if err != nil {
		logger.Fatalf("helpers:prompt", "failed to read input: %s", err.Error())
	}

	return result == "y"
}
