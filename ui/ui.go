package ui

import (
	"github.com/manifoldco/promptui"
)

func Prompt(label, defaultValue string) string {
	prompt := promptui.Prompt{
		Label:     label,
		Default:   defaultValue,
		AllowEdit: true,
	}

	input, err := prompt.Run()
	if err != nil {
		return ""
	}
	if input == "" {
		return ""
	}

	return input
}

func Confirm(label string) bool {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}

	_, err := prompt.Run()
	if err != nil {
		return false
	}

	return true
}

func Choose(label string, choices []string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: choices,
	}

	_, result, err := prompt.Run()

	return result, err
}
