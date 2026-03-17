package view

import (
	"fmt"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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

	form := m.forms.session.form
	if m.view == viewAdd {
		form = m.forms.add.form
	}
	if form != nil {
		b.WriteString(strings.TrimSuffix(form.View(), "\n"))
		b.WriteString("\n\n")
	}

	if m.status != "" {
		b.WriteString(
			join(helpStyle.Render(m.status), "\n\n"),
		)
	}

	b.WriteString(
		join(
			helpStyle.Render("? help"),
			helpStyle.Render("q quit"),
			"\n\n",
		),
	)

	return b.String()
}

func (m model) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	if key, ok := msg.(tea.KeyMsg); ok && m.view == viewAdd {
		switch key.String() {
		case "enter":
			return m.createSession(cmds)
		case "esc":
			m.view = viewList
			return m, tea.Batch(cmds...)
		}
	}

	if m.view == viewAdd {
		if m.forms.add.form != nil {
			form, cmd := m.forms.add.form.Update(msg)
			if f, ok := form.(*huh.Form); ok {
				m.forms.add.form = f
			}
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	}

	if m.forms.session.form != nil {
		form, cmd := m.forms.session.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.forms.session.form = f
		}
		if field := m.forms.session.form.GetFocusedField(); field != nil {
			if selectField, ok := field.(*huh.Select[string]); ok {
				if selected, ok := selectField.Hovered(); ok {
					m.forms.session.selected = selected
				}
				if selectField.GetFiltering() {
					return m, tea.Batch(cmds...)
				}
			}
		}
		cmds = append(cmds, cmd)
	}

	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, tea.Batch(cmds...)
	}

	switch key.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "?":
		m.view = viewHelp
		return m, tea.Batch(cmds...)

	case "c":
		m.view = viewAdd
		if m.forms.add.form != nil {
			cmds = append(cmds, m.forms.add.form.Init())
		}
		return m, tea.Batch(cmds...)

	case "enter":
		return m, tea.Batch(cmds...)

	case "d":
		if err := m.zellij.Delete(m.forms.session.selected); err != nil {
			m.status = err.Error()
			return m, tea.Batch(cmds...)
		}

		sessions, err := m.zellij.Ls()
		if err != nil {
			m.status = err.Error()
			return m, tea.Batch(cmds...)
		}
		m.status = "deleted: " + m.forms.session.selected
		m.forms.session.sessions = sessions
		m.view = viewList

		m.refreshSelectForm()
		if m.forms.session.form != nil {
			cmds = append(cmds, m.forms.session.form.Init())
		}

		return m, tea.Batch(cmds...)

	case "esc":
		if m.view == viewAdd {
			m.view = viewList
			return m, tea.Batch(cmds...)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) createSession(cmds []tea.Cmd) (tea.Model, tea.Cmd) {
	session := strings.TrimSpace(m.forms.add.value)
	if session == "" {
		m.status = "session name is empty"
		return m, tea.Batch(cmds...)
	}

	if err := m.zellij.Create(session); err != nil {
		m.status = err.Error()
		return m, tea.Batch(cmds...)
	}

	sessions, err := m.zellij.Ls()
	if err != nil {
		m.status = err.Error()
		return m, tea.Batch(cmds...)
	}

	m.forms.session.sessions = sessions
	m.forms.add.value = ""
	m.refreshAddForm()
	m.view = viewList
	m.refreshSelectForm()
	m.status = "created: " + session

	if m.forms.session.form != nil {
		cmds = append(cmds, m.forms.session.form.Init())
	}

	return m, tea.Batch(cmds...)
}

func (m *model) refreshAddForm() {
	m.forms.add.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("New session").
				Value(&m.forms.add.value),
		),
	).WithShowHelp(false).WithShowErrors(false)
}

func (m *model) refreshSelectForm() {
	var (
		options = make([]huh.Option[string], 0, len(m.forms.session.sessions))
	)

	for _, session := range m.forms.session.sessions {
		options = append(options, huh.NewOption(session, session))
	}

	if len(m.forms.session.sessions) > 0 {
		if m.forms.session.selected == "" {
			m.forms.session.selected = m.forms.session.sessions[0]
		} else {
			found := slices.Contains(m.forms.session.sessions, m.forms.session.selected)
			if !found {
				m.forms.session.selected = m.forms.session.sessions[0]
			}
		}
	}

	if len(options) == 0 {
		m.forms.session.selected = ""
		options = append(options, huh.NewOption("No sessions", ""))
	}

	m.forms.session.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose session").
				Description(fmt.Sprintf("Total: %d", len(m.forms.session.sessions))).
				Options(options...).
				Value(&m.forms.session.selected),
		),
	).WithShowHelp(false).WithShowErrors(false)
}
