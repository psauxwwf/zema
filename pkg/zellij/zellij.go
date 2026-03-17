package zellij

import (
	"fmt"
	"strings"
	"time"

	"zema/pkg/cmd"
)

type Zellij struct {
	binpath string
}

func New(
	_binpath string,
) (*Zellij, error) {
	return &Zellij{
		binpath: _binpath,
	}, nil
}

func (z *Zellij) cmd(args ...string) ([]byte, error) {
	_cmd, cancel, err := cmd.New(
		z.binpath,
		args,
		time.Second*3,
	)
	if err != nil {
		return nil, err
	}
	defer cancel()
	out, err := _cmd.Run()
	if err != nil {
		return out, fmt.Errorf("%w: %s", err, string(out))
	}
	return out, nil
}

func (z *Zellij) Ls() ([]string, error) {
	out, err := z.cmd("ls", "--short")
	if err != nil {
		return nil, fmt.Errorf("failed to ls sessions: %w", err)
	}
	return strings.Split(strings.TrimSuffix(string(out), "\n"), "\n"), nil
}

func (z *Zellij) Delete(name string) error {
	if _, err := z.cmd("delete-session", "--force", name); err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}
	return nil
}

func (z *Zellij) Create(name string) error {
	if _, err := z.cmd("attach", "--create-background", name); err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

func (z *Zellij) Attach(name string) error {
	if _, err := z.cmd("attach", "--create", name); err != nil {
		return fmt.Errorf("failed attach to session: %w", err)
	}
	return nil
}
