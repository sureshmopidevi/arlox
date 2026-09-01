package cli

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"golang.org/x/term"

	"github.com/sureshmopidevi/arlox/internal/designsystems"
	"github.com/sureshmopidevi/arlox/internal/workspace"
)

func stacksIncludeWeb(stacks []workspace.Stack) bool {
	for _, s := range stacks {
		if s == workspace.Web {
			return true
		}
	}
	return false
}

func resolveWebDesignSystem(f stackFlags, stacks []workspace.Stack) (string, error) {
	if !stacksIncludeWeb(stacks) {
		return "", nil
	}
	if f.webUI != "" {
		if err := designsystems.ValidateWebID(f.webUI); err != nil {
			return "", err
		}
		return f.webUI, nil
	}
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return designsystems.DefaultWebID, nil
	}
	return promptWebDesignSystem()
}

func promptWebDesignSystem() (string, error) {
	opts := make([]huh.Option[string], len(designsystems.WebCatalog))
	for i, s := range designsystems.WebCatalog {
		label := s.Label
		if s.Description != "" {
			label = s.Label + " — " + s.Description
		}
		opts[i] = huh.NewOption(label, s.ID)
	}

	var selected string
	err := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().
			Title("Web design system").
			Description("UI library for the React admin console").
			Options(opts...).
			Value(&selected),
	)).Run()
	if err != nil {
		return "", err
	}
	if selected == "" {
		return "", fmt.Errorf("no design system selected")
	}
	return selected, nil
}
