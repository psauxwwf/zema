package view

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	helpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	okStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
	badStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	waitStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("246"))
	selectStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
)
