package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"

	"zema/internal/config"
	"zema/internal/zema"
)

const (
	_ int = iota
	defaultCode
	configCode
	initCode
	fatalCode
)

type exitError struct {
	code int
	err  error
}

func (e *exitError) Error() string {
	if e == nil || e.err == nil {
		return ""
	}
	return e.err.Error()
}

func (e *exitError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

func newExitError(code int, err error) error {
	if err == nil {
		return nil
	}

	return &exitError{code: code, err: err}
}

func main() {
	if err := fang.Execute(context.Background(), rootCmd(), fang.WithoutVersion()); err != nil {
		if err, ok := errors.AsType[*exitError](err); ok {
			fmt.Println(err.err)
			os.Exit(err.code)
		}
		fmt.Println(err)
		os.Exit(fatalCode)
	}
}

func rootCmd() *cobra.Command {
	var (
		path        string
		save        bool
		completion  bool
		tabTitlePwd string
		tabTitleCmd string
	)

	root := &cobra.Command{
		Use:           "zema [session]",
		Short:         "Terminal UI for managing Zellij sessions",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			path = toAbs(path)

			if completion {
				fmt.Print(zshCompletionScript())
				return nil
			}

			if save {
				if err := config.Default(path); err != nil {
					return newExitError(defaultCode, err)
				}
				return nil
			}

			_config, err := config.New(path)
			if err != nil && !errors.Is(err, config.ErrNotExists) {
				return newExitError(configCode, err)
			}

			_zema, err := zema.New(_config)
			if err != nil {
				return newExitError(initCode, err)
			}

			if len(args) > 0 {
				if err := _zema.Attach(args[0]); err != nil {
					return newExitError(fatalCode, err)
				}
				return nil
			}

			if tabTitlePwd != "" || tabTitleCmd != "" {
				if tabTitlePwd != "" {
					_ = _zema.RenameTabFromPwd(tabTitlePwd)
					return nil
				}
				if tabTitleCmd != "" {
					_ = _zema.RenameTabFromCommand(tabTitleCmd)
					return nil
				}
			}

			if err := _zema.Run(); err != nil {
				return newExitError(fatalCode, err)
			}

			return nil
		},
	}

	root.Flags().StringVar(&path, "config", "~/.config/zema/config.yaml", "path to config")
	root.Flags().BoolVar(&save, "save", false, "save default config to path --config")
	root.Flags().BoolVar(&completion, "completion", false, "print zsh completion script for eval")
	root.Flags().StringVar(&tabTitlePwd, "tab-title-pwd", "", "set zellij tab title from current working directory")
	root.Flags().StringVar(&tabTitleCmd, "tab-title-cmd", "", "set zellij tab title from executed command")

	return root
}

func toAbs(path string) string {
	if path == "~" {
		home, err := os.UserHomeDir()
		if err == nil {
			return home
		}
		return path
	}

	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		path = filepath.Join(home, strings.TrimPrefix(path, "~/"))
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}

	return abs
}

func zshCompletionScript() string {
	return `func _zema_tab_precmd() { zema --tab-title-pwd "$PWD" &>/dev/null }
func _zema_tab_preexec() { zema --tab-title-cmd "$1" &>/dev/null }
autoload -Uz add-zsh-hook
add-zsh-hook precmd _zema_tab_precmd
add-zsh-hook preexec _zema_tab_preexec
add-zsh-hook chpwd _zema_tab_precmd
`
}
