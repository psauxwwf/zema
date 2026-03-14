package view

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

func (m model) viewList() string {
	var b strings.Builder

	b.WriteString(
		join(
			titleStyle.Render("zema"),
			helpStyle.Render("zellij manager"),
			"\n\n",
		),
	)
	b.WriteString(
		join(
			helpStyle.Render("? help"),
			helpStyle.Render("q quit"),
			"\n\n",
		),
	)

	switch {
	case m.filter.Value() != "":
		{
			b.WriteString(
				join(">", m.filter.Value(), "\n\n"),
			)
			break
		}
	case m.isFilter:
		{
			b.WriteString(
				join(m.filter.View(), "\n\n"),
			)

		}
	}

	if len(m.filtered) == 0 {
		b.WriteString(
			join(helpStyle.Render("No matches"), "\n"),
		)
	} else {
		for fi, idx := range m.filtered {
			id := m.sessions[idx]
			line := "  " + id
			if fi == m.cursor {
				line = selectStyle.Render("> " + id)
			}
			b.WriteString(line + "\n")
		}
	}

	if m.selectForm != nil {
		b.WriteString("\n")
		b.WriteString(strings.TrimSuffix(m.selectForm.View(), "\n"))
	}

	if m.status != "" {
		b.WriteString(
			join("\n", helpStyle.Render(m.status), "\n"),
		)
	}

	return b.String()
}

func (m model) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if m.selectForm != nil {
		form, cmd := m.selectForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.selectForm = f
		}
		cmds = append(cmds, cmd)
	}

	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, tea.Batch(cmds...)
	}

	m.status = ""

	if m.isFilter {
		switch key.String() {
		case "enter":
			m.isFilter = false
			return m, tea.Batch(cmds...)
		case "esc":
			m.isFilter = false
			m.resetFilter()
			return m, tea.Batch(cmds...)
		}

		var cmd tea.Cmd
		m.filter, cmd = m.filter.Update(msg)
		m.applyFilter(m.filter.Value())
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	switch key.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "?":
		m.view = viewHelp
		return m, tea.Batch(cmds...)
	case "/":
		m.isFilter = true
		m.filter.Focus()
		cmds = append(cmds, textinput.Blink)
		return m, tea.Batch(cmds...)
	case "esc":
		m.filter.SetValue("")
		m.resetFilter()
		return m, tea.Batch(cmds...)
	case "a":
		m.view = viewAdd
		m.add.Focus()
		cmds = append(cmds, textinput.Blink)
		return m, tea.Batch(cmds...)
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.filtered)-1 {
			m.cursor++
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *model) resetFilter() {
	m.filter.SetValue("")
	m.refreshSelectForm()
}

func (m *model) applyFilter(query string) {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		m.resetFilter()
		return
	}

	m.cursor = 0
	m.refreshSelectForm()
}

func (m *model) refreshSelectForm() {
	options := make([]huh.Option[string], 0, len(m.sessions))
	for _, id := range m.sessions {
		options = append(options, huh.NewOption(id, id))
	}

	if len(options) == 0 {
		options = append(options, huh.NewOption("No sessions", ""))
	}

	m.selectForm = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose session").
				Description(fmt.Sprintf("Total: %d", len(m.sessions))).
				Options(options...).
				Value(&m.selectedID),
		),
	).WithShowHelp(false).WithShowErrors(false)
}
