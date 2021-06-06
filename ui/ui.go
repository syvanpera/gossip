package ui

import (
	"io/ioutil"
	"os"
	"os/exec"

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

func Editor(content string) string {
	f, err := ioutil.TempFile("", "gossip")
	if err != nil {
		return ""
	}
	defer os.Remove(f.Name())

	if _, err := f.WriteString(content); err != nil {
		return ""
	}

	f.Close()

	editor, _ := exec.LookPath(os.Getenv("EDITOR"))
	command := exec.Command(editor, f.Name())

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err = command.Start(); err != nil {
		return ""
	}
	if err = command.Wait(); err != nil {
		return ""
	}

	data, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return ""
	}
	return string(data)
}
