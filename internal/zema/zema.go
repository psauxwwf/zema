package zema

import (
	"fmt"
	"zema/internal/view"
	"zema/pkg/zellij"

	tea "github.com/charmbracelet/bubbletea"
)

type Zema struct {
	view *tea.Program
}

func New() (*Zema, error) {
	_zellij, err := zellij.New("zellij")
	if err != nil {
		return nil, fmt.Errorf("fatal zellij: %w", err)
	}
	_view, err := view.New(
		_zellij,
	)
	if err != nil {
		return nil, fmt.Errorf("fatal interface: %w", err)
	}
	return &Zema{
		view: _view,
	}, nil
}

func (z *Zema) Run() error {
	if _, err := z.view.Run(); err != nil {
		return err
	}
	return nil
}
