package view

import (
	"os/exec"
	"strings"

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
	Delete(name string) error
	Create(name string) error
	Attach(name string) *exec.Cmd
}

type size struct {
	width, height int
}

type forms struct {
	sessions *sessions
	add      *add
}

type sessions struct {
	form     *huh.Form
	selected string
	sessions []string
}

type add struct {
	form  *huh.Form
	value string
}

type model struct {
	zellij Zellij
	view   view
	status string

	size  *size
	forms *forms
}

func New(
	_zellij Zellij,
) (*tea.Program, error) {
	var (
		_status        string
		_sessions, err = _zellij.Ls()
	)
	if err != nil {
		_status = strings.TrimSpace(err.Error())
	}

	m := model{
		zellij: _zellij,
		size:   &size{},
		forms: &forms{
			sessions: &sessions{
				sessions: _sessions,
			},
			add: &add{},
		},
		status: _status,
	}
	m.refreshAddForm()
	m.refreshSessionsForm()

	return tea.NewProgram(m, tea.WithAltScreen()), nil
}

func (m model) Init() tea.Cmd {
	if m.forms.sessions.form != nil {
		return tea.Batch(m.forms.sessions.form.Init())
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
	return strings.Join(s, joinSeparator)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size.width = msg.Width
		m.size.height = msg.Height
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
