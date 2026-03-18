package zellij

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"zema/internal/config"
	"zema/pkg/cmd"
)

type Zellij struct {
	binpath string
	ls      []string
	delete  []string
	create  []string
	attach  config.Command
}

func New(
	config config.Zellij,
) (*Zellij, error) {
	return &Zellij{
		binpath: config.Bin,
		ls:      config.Ls,
		delete:  config.Delete,
		create:  config.Create,
		attach:  config.Attach,
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
	out, err := z.cmd(z.ls...)
	if err != nil {
		return []string{}, fmt.Errorf("failed to ls sessions: %w", err)
	}
	return strings.Split(strings.TrimSuffix(string(out), "\n"), "\n"), nil
}

func (z *Zellij) Delete(name string) error {
	if _, err := z.cmd(renderArgs(z.delete, name, "")...); err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}
	return nil
}

func (z *Zellij) Create(name string) error {
	if _, err := z.cmd(renderArgs(z.create, name, "")...); err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

func (z *Zellij) Attach(name string) *exec.Cmd {
	home, err := os.UserHomeDir()
	if err != nil {
		home = ""
	}

	var (
		args   = renderArgs(z.attach.Args, name, home)
		pre    = shellJoin(renderArgs(z.attach.Pre, name, home))
		attach = shellJoin(append([]string{z.binpath}, args...))
		post   = shellJoin(renderArgs(z.attach.Post, name, home))
	)

	if pre == "" || post == "" {
		return exec.Command(z.binpath, args...)
	}
	script := fmt.Sprintf("%s && %s && %s", pre, attach, post)

	return exec.Command("sh", "-c", script)
}

func renderArgs(args []string, name string, home string) []string {
	var (
		prepared       = make([]string, len(args))
		hasPlaceholder = false
	)

	for i, arg := range args {
		if arg == "{session}" {
			prepared[i] = name
			hasPlaceholder = true
			continue
		}
		rendered := strings.ReplaceAll(arg, "{session}", name)
		rendered = strings.ReplaceAll(rendered, "{home}", home)
		if strings.HasPrefix(rendered, "~/") && home != "" {
			rendered = strings.Replace(rendered, "~", home, 1)
		}
		prepared[i] = rendered
	}

	if hasPlaceholder {
		return prepared
	}

	return append(prepared, name)
}

func shellJoin(args []string) string {
	quoted := make([]string, len(args))
	for i, arg := range args {
		quoted[i] = shellQuote(arg)
	}
	return strings.Join(quoted, " ")
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
}
