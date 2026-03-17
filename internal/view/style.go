package view

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#BD93F9"))
	subtitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#8BE9FD"))
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#6272A4"))
	okStyle       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50FA7B"))
	badStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF5555"))
	waitStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB86C"))
	hotkeyStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F1FA8C"))
	panelStyle    = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#44475A")).
			Padding(0, 1)
)

func renderStatus(status string) string {
	value := strings.TrimSpace(status)
	if value == "" {
		return ""
	}

	s := strings.ToLower(value)
	switch {
	case strings.HasPrefix(s, statusCreatedPrefix), strings.HasPrefix(s, statusDeletedPrefix):
		return okStyle.Render(value)
	case strings.Contains(s, "failed"), strings.Contains(s, "fatal"), strings.Contains(s, "empty"), strings.Contains(s, "not selected"):
		return badStyle.Render(value)
	default:
		return waitStyle.Render(value)
	}
}

func keyHint(key, action string) string {
	return hotkeyStyle.Render(key) + " " + helpStyle.Render(action)
}
