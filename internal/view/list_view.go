package view

import (
	"strings"

	"github.com/charmbracelet/huh"
)

func (m model) viewList() string {
	var b strings.Builder

	b.WriteString(
		join(
			titleStyle.Render("Zema:"),
			helpStyle.Render("zellij manager"),
			"\n\n",
		),
	)

	if form := m.currForm(); form != nil {
		b.WriteString(strings.TrimSuffix(form.View(), "\n"))
		b.WriteString("\n\n")
	}

	if m.status != "" {
		b.WriteString(join(helpStyle.Render(m.status), "\n\n"))
	}

	b.WriteString(
		join(
			helpStyle.Render("q quit"),
			helpStyle.Render("? help"),
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
