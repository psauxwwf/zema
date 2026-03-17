package view

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type attachDoneMsg struct {
	err error
}

func (m model) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	if done, ok := msg.(attachDoneMsg); ok {
		return m.onDone(done)
	}

	if key, ok := msg.(tea.KeyMsg); ok && m.view == viewAdd {
		return m.onAdd(key)
	}

	if m.view == viewAdd {
		return m.updateAddForm(msg)
	}

	if key, ok := msg.(tea.KeyMsg); ok {
		next, cmd, handled := m.onKey(key)
		if handled {
			return next, cmd
		}
	}

	return m.updateSessionForm(msg)
}

func (m model) onDone(done attachDoneMsg) (tea.Model, tea.Cmd) {
	if done.err != nil {
		m.status = strings.TrimSpace(done.err.Error())
		return m, nil
	}

	return m, tea.Quit
}

func (m model) onAdd(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "enter":
		return m.addSession()
	case "esc":
		m.view = viewList
		return m, nil
	default:
		return m.updateAddForm(key)
	}
}

func (m model) updateAddForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.forms.add.form == nil {
		return m, nil
	}

	form, cmd := m.forms.add.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.forms.add.form = f
	}

	return m, cmd
}

func (m model) onKey(key tea.KeyMsg) (tea.Model, tea.Cmd, bool) {
	switch key.String() {
	case "q", "ctrl+c":
		return m, tea.Quit, true
	case "?":
		m.view = viewHelp
		return m, nil, true
	case "c":
		m.view = viewAdd
		return m, initForm(m.forms.add.form), true
	case "enter":
		next, cmd := m.attachSession()
		return next, cmd, true
	case "d":
		next, cmd := m.deleteSession()
		return next, cmd, true
	default:
		return m, nil, false
	}
}

func (m model) attachSession() (tea.Model, tea.Cmd) {
	if m.forms.sessions.selected == "" {
		m.status = "session is not selected"
		return m, nil
	}

	return m, tea.ExecProcess(
		m.zellij.Attach(m.forms.sessions.selected),
		func(err error) tea.Msg {
			return attachDoneMsg{err: err}
		},
	)
}

func (m model) deleteSession() (tea.Model, tea.Cmd) {
	if err := m.zellij.Delete(m.forms.sessions.selected); err != nil {
		m.status = strings.TrimSpace(err.Error())
		return m, nil
	}

	sessions, err := m.zellij.Ls()
	if err != nil {
		m.status = strings.TrimSpace(err.Error())
		return m, nil
	}

	m.status = "deleted: " + m.forms.sessions.selected
	m.forms.sessions.sessions = sessions
	m.view = viewList
	m.refreshSessionsForm()

	return m, initForm(m.forms.sessions.form)
}

func (m model) updateSessionForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.forms.sessions.form == nil {
		return m, nil
	}

	form, cmd := m.forms.sessions.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.forms.sessions.form = f
	}

	if field := m.forms.sessions.form.GetFocusedField(); field != nil {
		if sf, ok := field.(*huh.Select[string]); ok {
			if selected, ok := sf.Hovered(); ok {
				m.forms.sessions.selected = selected
			}
			if sf.GetFiltering() {
				return m, nil
			}
		}
	}

	return m, cmd
}

func (m model) addSession() (tea.Model, tea.Cmd) {
	session := strings.TrimSpace(m.forms.add.value)
	if session == "" {
		m.status = "session name is empty"
		return m, nil
	}

	if err := m.zellij.Create(session); err != nil {
		m.status = strings.TrimSpace(err.Error())
		return m, nil
	}

	sessions, err := m.zellij.Ls()
	if err != nil {
		m.status = strings.TrimSpace(err.Error())
		return m, nil
	}

	m.forms.sessions.sessions = sessions
	m.forms.add.value = ""
	m.refreshAddForm()
	m.view = viewList
	m.refreshSessionsForm()
	m.status = "created: " + session

	return m, initForm(m.forms.sessions.form)
}
