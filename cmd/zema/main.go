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
	path = flag.String("config", "~/.config/zema/config.yaml", "path to config")
	save = flag.Bool("save", false, "save default config to path -config")
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
