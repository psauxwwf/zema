package view

import (
	"strings"

	"github.com/charmbracelet/huh"
)

func (m model) viewList() string {
	var b strings.Builder

	header := join(
		titleStyle.Render(appTitle),
		subtitleStyle.Render(sessionsTitle),
	)

	b.WriteString(
		panelStyle.Render(header),
	)
	b.WriteString("\n")

	if form := m.currForm(); form != nil {
		b.WriteString(strings.TrimSuffix(form.View(), "\n"))
		b.WriteString("\n\n")
	}

	if m.status != "" {
		b.WriteString(renderStatus(m.status))
		b.WriteString("\n\n")
	}

	b.WriteString(
		join(
			keyHint(keyQuit, actionQuit),
			keyHint(keyHelp, actionHelp),
			"\n\n",
		),
	)

	return b.String()
}

func (m model) currForm() *huh.Form {
	if m.view == viewAdd {
		return m.forms.add.form
	}
	return m.forms.sessions.form
}
