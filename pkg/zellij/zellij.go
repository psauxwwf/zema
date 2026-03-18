package zellij

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"zema/internal/config"
	"zema/pkg/cmd"
)

type Zellij struct {
	binpath   string
	ls        []string
	delete    []string
	create    []string
	renameTab []string
	attach    config.Command
}

func New(
	config config.Zellij,
) (*Zellij, error) {
	return &Zellij{
		binpath:   config.Bin,
		ls:        config.Ls,
		delete:    config.Delete,
		create:    config.Create,
		renameTab: config.RenameTab,
		attach:    config.Attach,
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
	if _, err := z.cmd(renderArgs(z.delete, name, "", "{session}")...); err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}
	return nil
}

func (z *Zellij) Create(name string) error {
	if _, err := z.cmd(renderArgs(z.create, name, "", "{session}")...); err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

func (z *Zellij) RenameTab(title string) error {
	if os.Getenv("ZELLIJ") == "" {
		return nil
	}

	title = strings.TrimSpace(title)
	if title == "" {
		return nil
	}

	runes := []rune(title)
	if len(runes) > 32 {
		title = string(runes[:32])
	}

	if _, err := z.cmd(renderArgs(z.renameTab, title, "", "{title}")...); err != nil {
		return fmt.Errorf("failed to rename tab: %w", err)
	}

	return nil
}

func (z *Zellij) RenameTabFromPwd(pwd string) error {
	home, err := os.UserHomeDir()
	if err == nil && pwd == home {
		return z.RenameTab("~")
	}

	return z.RenameTab(filepath.Base(pwd))
}

func (z *Zellij) RenameTabFromCommand(command string) error {
	return z.RenameTab(strings.TrimSpace(command))
}

func (z *Zellij) Attach(name string) *exec.Cmd {
	home, err := os.UserHomeDir()
	if err != nil {
		home = ""
	}

	var (
		args   = renderArgs(z.attach.Args, name, home, "{session}")
		pre    = shellJoin(renderArgs(z.attach.Pre, name, home, "{session}"))
		attach = shellJoin(append([]string{z.binpath}, args...))
		post   = shellJoin(renderArgs(z.attach.Post, name, home, "{session}"))
	)

	if pre == "" || post == "" {
		return exec.Command(z.binpath, args...)
	}
	script := fmt.Sprintf("%s && %s && %s", pre, attach, post)

	return exec.Command("sh", "-c", script)
}

func renderArgs(args []string, value string, home string, placeholder string) []string {
	var (
		prepared       = make([]string, len(args))
		hasPlaceholder = false
	)

	for i, arg := range args {
		if arg == placeholder {
			prepared[i] = value
			hasPlaceholder = true
			continue
		}
		rendered := strings.ReplaceAll(arg, placeholder, value)
		rendered = strings.ReplaceAll(rendered, "{home}", home)
		if strings.HasPrefix(rendered, "~/") && home != "" {
			rendered = strings.Replace(rendered, "~", home, 1)
		}
		prepared[i] = rendered
	}

	if hasPlaceholder {
		return prepared
	}

	return append(prepared, value)
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
