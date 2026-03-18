package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"zema/internal/config"
	"zema/internal/zema"
)

type code int

var (
	path        = flag.String("config", "~/.config/zema/config.yaml", "path to config")
	save        = flag.Bool("save", false, "save default config to path -config")
	completion  = flag.Bool("completion", false, "print zsh completion script for eval")
	tabTitlePwd = flag.String("tab-title-pwd", "", "set zellij tab title from current working directory")
	tabTitleCmd = flag.String("tab-title-cmd", "", "set zellij tab title from executed command")
)

const (
	_ int = iota
	defaultCode
	configCode
	initCode
	fatalCode
)

func main() {
	flag.Parse()

	if *completion {
		fmt.Print(zshCompletionScript())
		return
	}

	path := toAbs(*path)

	if *save {
		if err := config.Default(path); err != nil {
			fmt.Println(err)
			os.Exit(defaultCode)
		}
		return
	}

	_config, err := config.New(path)
	if err != nil && !errors.Is(err, config.ErrNotExists) {
		fmt.Println(err)
		os.Exit(configCode)
	}

	_zema, err := zema.New(_config)
	if err != nil {
		fmt.Println(err)
		os.Exit(initCode)
	}

	if flag.NArg() > 0 {
		if err := _zema.Attach(flag.Arg(0)); err != nil {
			fmt.Println(err)
			os.Exit(fatalCode)
		}
		return
	}

	if *tabTitlePwd != "" || *tabTitleCmd != "" {
		if *tabTitlePwd != "" {
			_ = _zema.RenameTabFromPwd(*tabTitlePwd)
			return
		}
		if *tabTitleCmd != "" {
			_ = _zema.RenameTabFromCommand(*tabTitleCmd)
			return
		}
	}

	if err := _zema.Run(); err != nil {
		fmt.Println(err)
		os.Exit(fatalCode)
	}
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
	return `func _zema_tab_precmd() { zema -tab-title-pwd "$PWD" &>/dev/null }
func _zema_tab_preexec() { zema -tab-title-cmd "$1" &>/dev/null }
autoload -Uz add-zsh-hook
add-zsh-hook precmd _zema_tab_precmd
add-zsh-hook preexec _zema_tab_preexec
`
}
