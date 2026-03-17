package view

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const helpKeyColumnWidth = 16

func helpHint(key, action string) string {
	keyCol := lipgloss.NewStyle().Width(helpKeyColumnWidth).Render(hotkeyStyle.Render(key))
	return keyCol + helpStyle.Render(action)
}

func (m model) viewHelp() string {
	return strings.Join([]string{
		panelStyle.Render(join(titleStyle.Render("Zema 🧌"), subtitleStyle.Render("help"))),
		helpHint("j/k or up/down", "move select"),
		helpHint("enter", "attach to session"),
		helpHint("/", "filter sessions"),
		helpHint("c", "create session"),
		helpHint("d", "delete session"),
		helpHint("q", "quit"),
		"",
		waitStyle.Render("Press any key to close help"),
	}, "\n")
}
