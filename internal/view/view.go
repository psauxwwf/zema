package view

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type view int

const (
	viewList view = iota
	viewHelp
	viewAdd
)

type Zellij interface {
	Ls() ([]string, error)
}

type model struct {
	zellij Zellij

	width, height int

	view view

	sessions []string

	filter, add textinput.Model
	selectForm  *huh.Form
	selectedID  string

	isFilter bool
	filtered []int

	cursor int
	status string
}

func New(
	_zellij Zellij,
) (*tea.Program, error) {
	_filter := textinput.New()
	_filter.Placeholder = "Search session..."
	_filter.CharLimit = 64
	_filter.Width = 30

	_add := textinput.New()
	_add.Placeholder = "New session"
	_add.CharLimit = 64
	_add.Width = 30

	_sessions, err := _zellij.Ls()
	if err != nil {
		return nil, err
	}

	m := model{
		zellij:   _zellij,
		filter:   _filter,
		add:      _add,
		sessions: _sessions,
	}
	m.refreshSelectForm()

	return tea.NewProgram(m, tea.WithAltScreen()), nil
}

func (m model) Init() tea.Cmd {
	if m.selectForm != nil {
		return tea.Batch(m.selectForm.Init())
	}
	return tea.Batch()
}

func (m model) View() string {
	switch m.view {
	case viewHelp:
		return m.viewHelp()
	default:
		return m.viewList()
	}
}

func join(s ...string) string {
	return strings.Join(s, "  ")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	switch m.view {
	case viewHelp:
		if _, ok := msg.(tea.KeyMsg); ok {
			m.view = viewList
		}
		return m, nil
	default:
		return m.updateList(msg)
	}
}
