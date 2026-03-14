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
	return _cmd.Run()
}

func (z *Zellij) Ls() ([]string, error) {
	out, err := z.cmd("ls", "--short")
	if err != nil {
		return nil, fmt.Errorf("ls error: %w", err)
	}

	return strings.Split(strings.TrimSuffix(string(out), "\n"), "\n"), nil
}
