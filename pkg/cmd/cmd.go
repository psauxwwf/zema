package cmd

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"time"
)

type Cmd struct {
	ctx    context.Context
	cmd    *exec.Cmd
	cancel context.CancelFunc
}

var (
	ErrTimeout = errors.New("timeout")
)

func New(
	command string,
	args []string,
	timeout time.Duration,
	envs ...string,
) (*Cmd, context.CancelFunc, error) {
	var (
		_ctx, _cancel = context.WithTimeout(context.Background(), timeout)
		_cmd          = exec.CommandContext(_ctx, command, args...)
	)
	_cmd.Env = append(os.Environ(),
		envs...,
	)
	return &Cmd{
		ctx:    _ctx,
		cmd:    _cmd,
		cancel: _cancel,
	}, _cancel, nil
}

func (cmd *Cmd) Run() ([]byte, error) {
	out, err := cmd.cmd.CombinedOutput()
	if err != nil {
		if errors.Is(cmd.ctx.Err(), context.DeadlineExceeded) {
			return out, ErrTimeout
		}
		return out, err
	}
	return out, nil
}
