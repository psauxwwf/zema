package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type viewState int

const (
	viewList viewState = iota
	viewHelp
	viewAdd
)

type checkResultMsg struct {
	id string
	ok bool
}

type model struct {
	items   []string
	health  map[string]string
	view    viewState
	cursor  int
	width   int
	height  int
	status  string
	version string

	filtering bool
	filter    textinput.Model
	filtered  []int

	addInput textinput.Model
}

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	okStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
	badStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	waitStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("246"))
	selStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
)

func newModel() model {
	f := textinput.New()
	f.Placeholder = "Type to filter..."
	f.CharLimit = 64
	f.Width = 30

	a := textinput.New()
	a.Placeholder = "new-server"
	a.CharLimit = 64
	a.Width = 30

	m := model{
		items: []string{
			"dev-web-1",
			"dev-db-1",
			"staging-api",
			"prod-web-1",
			"prod-db-1",
		},
		health:   map[string]string{},
		view:     viewList,
		version:  "0.1.0",
		filter:   f,
		addInput: a,
	}

	m.resetFilter()
	for _, id := range m.items {
		m.health[id] = "checking"
	}

	return m
}

func (m model) Init() tea.Cmd {
	return m.buildHealthChecksCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case checkResultMsg:
		if msg.ok {
			m.health[msg.id] = "reachable"
		} else {
			m.health[msg.id] = "unreachable"
		}
		return m, nil
	}

	switch m.view {
	case viewHelp:
		if _, ok := msg.(tea.KeyMsg); ok {
			m.view = viewList
		}
		return m, nil
	case viewAdd:
		return m.updateAdd(msg)
	default:
		return m.updateList(msg)
	}
}

func (m model) updateAdd(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "esc":
			m.view = viewList
			m.addInput.SetValue("")
			return m, nil
		case "enter":
			id := strings.TrimSpace(m.addInput.Value())
			if id == "" {
				m.status = "ID is empty"
				m.view = viewList
				return m, nil
			}
			m.items = append(m.items, id)
			m.health[id] = "checking"
			m.addInput.SetValue("")
			m.view = viewList
			m.status = fmt.Sprintf("Added: %s", id)
			m.resetFilter()
			return m, checkHealthCmd(id)
		}
	}

	var cmd tea.Cmd
	m.addInput, cmd = m.addInput.Update(msg)
	return m, cmd
}

func (m model) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	m.status = ""

	if m.filtering {
		switch key.String() {
		case "enter":
			m.filtering = false
			return m, nil
		case "esc":
			m.filtering = false
			m.filter.SetValue("")
			m.resetFilter()
			return m, nil
		}

		var cmd tea.Cmd
		m.filter, cmd = m.filter.Update(msg)
		m.applyFilter(m.filter.Value())
		return m, cmd
	}

	switch key.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "?":
		m.view = viewHelp
		return m, nil
	case "/":
		m.filtering = true
		m.filter.Focus()
		return m, textinput.Blink
	case "esc":
		m.filter.SetValue("")
		m.resetFilter()
		return m, nil
	case "a":
		m.view = viewAdd
		m.addInput.Focus()
		return m, textinput.Blink
	case "r":
		for _, id := range m.items {
			m.health[id] = "checking"
		}
		m.status = "Refreshing health checks..."
		return m, m.buildHealthChecksCmd()
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.filtered)-1 {
			m.cursor++
		}
	}

	return m, nil
}

func (m *model) resetFilter() {
	m.filtered = m.filtered[:0]
	for i := range m.items {
		m.filtered = append(m.filtered, i)
	}
	if m.cursor >= len(m.filtered) {
		m.cursor = max(0, len(m.filtered)-1)
	}
}

func (m *model) applyFilter(query string) {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		m.resetFilter()
		return
	}

	m.filtered = m.filtered[:0]
	for i, id := range m.items {
		if strings.Contains(strings.ToLower(id), query) {
			m.filtered = append(m.filtered, i)
		}
	}
	m.cursor = 0
}

func (m model) buildHealthChecksCmd() tea.Cmd {
	cmds := make([]tea.Cmd, 0, len(m.items))
	for _, id := range m.items {
		cmds = append(cmds, checkHealthCmd(id))
	}
	return tea.Batch(cmds...)
}

func checkHealthCmd(id string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(250 * time.Millisecond)
		ok := !strings.Contains(strings.ToLower(id), "prod-db")
		return checkResultMsg{id: id, ok: ok}
	}
}

func (m model) View() string {
	switch m.view {
	case viewHelp:
		return m.viewHelp()
	case viewAdd:
		return m.viewAdd()
	default:
		return m.viewList()
	}
}

func (m model) viewHelp() string {
	return strings.Join([]string{
		titleStyle.Render("Help"),
		"",
		helpStyle.Render("j/k or up/down  move"),
		helpStyle.Render("/               filter"),
		helpStyle.Render("a               add item"),
		helpStyle.Render("r               refresh checks"),
		helpStyle.Render("q               quit"),
		"",
		helpStyle.Render("Press any key to close help"),
	}, "\n")
}

func (m model) viewAdd() string {
	return strings.Join([]string{
		titleStyle.Render("Add Connection"),
		"",
		"ID: " + m.addInput.View(),
		"",
		helpStyle.Render("enter save   esc cancel"),
	}, "\n")
}

func (m model) viewList() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("zema") + helpStyle.Render("  Bubble Tea demo") + "\n")
	b.WriteString(helpStyle.Render("/ filter  a add  r refresh  ? help  q quit") + "\n\n")

	if m.filtering {
		b.WriteString("Filter: " + m.filter.View() + "\n\n")
	} else if m.filter.Value() != "" {
		b.WriteString("Filter: " + m.filter.Value() + "\n\n")
	} else {
		b.WriteString(helpStyle.Render(fmt.Sprintf("%d items", len(m.filtered))) + "\n\n")
	}

	if len(m.filtered) == 0 {
		b.WriteString(helpStyle.Render("No matches") + "\n")
	} else {
		for fi, idx := range m.filtered {
			id := m.items[idx]
			dot := renderHealth(m.health[id])
			line := "  " + id + " " + dot
			if fi == m.cursor {
				line = selStyle.Render("> "+id) + " " + dot
			}
			b.WriteString(line + "\n")
		}
	}

	if m.status != "" {
		b.WriteString("\n" + helpStyle.Render(m.status) + "\n")
	}

	return b.String()
}

func renderHealth(status string) string {
	switch status {
	case "reachable":
		return okStyle.Render("●")
	case "unreachable":
		return badStyle.Render("●")
	default:
		return waitStyle.Render("○")
	}
}

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("TUI error:", err)
	}
}
