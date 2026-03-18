package zema

import (
	"fmt"
	"os"
	"zema/internal/config"
	"zema/internal/view"
	"zema/pkg/zellij"

	tea "github.com/charmbracelet/bubbletea"
)

type Zema struct {
	view   *tea.Program
	zellij *zellij.Zellij
}

func New(config *config.Config) (*Zema, error) {
	_zellij, err := zellij.New(config.Zellij)
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
		view:   _view,
		zellij: _zellij,
	}, nil
}

func (z *Zema) Run() error {
	if _, err := z.view.Run(); err != nil {
		return err
	}
	return nil
}

func (z *Zema) RenameTabFromPwd(pwd string) error {
	return z.zellij.RenameTabFromPwd(pwd)
}

func (z *Zema) RenameTabFromCommand(command string) error {
	return z.zellij.RenameTabFromCommand(command)
}

func (z *Zema) Attach(name string) error {
	cmd := z.zellij.Attach(name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to attach session: %w", err)
	}

	return nil
}
