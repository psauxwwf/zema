package view

import (
	"fmt"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

func initForm(form *huh.Form) tea.Cmd {
	if form == nil {
		return nil
	}
	return form.Init()
}

func (m *model) refreshAddForm() {
	m.forms.add.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(labelNewSession).
				Value(&m.forms.add.value),
		),
	).WithTheme(huh.ThemeDracula()).WithShowHelp(false).WithShowErrors(false)
}

func (m *model) refreshSessionsForm() {
	options := make([]huh.Option[string], 0, len(m.forms.sessions.sessions))

	for _, session := range m.forms.sessions.sessions {
		options = append(options, huh.NewOption(session, session))
	}

	m.fixSelect()

	if len(options) == 0 {
		m.forms.sessions.selected = ""
		options = append(options, huh.NewOption(labelNoSessions, ""))
	}

	m.forms.sessions.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(labelChooseSession).
				Description(fmt.Sprintf(descTotalFmt, len(m.forms.sessions.sessions))).
				Options(options...).
				Value(&m.forms.sessions.selected),
		),
	).WithTheme(huh.ThemeDracula()).WithShowHelp(false).WithShowErrors(false)
}

func (m *model) fixSelect() {
	if len(m.forms.sessions.sessions) == 0 {
		return
	}

	if m.forms.sessions.selected == "" {
		m.forms.sessions.selected = m.forms.sessions.sessions[0]
		return
	}

	if !slices.Contains(m.forms.sessions.sessions, m.forms.sessions.selected) {
		m.forms.sessions.selected = m.forms.sessions.sessions[0]
	}
}
