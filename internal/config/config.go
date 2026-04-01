package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var (
	_default = Config{
		Zellij: Zellij{
			Bin:       "zellij",
			Ls:        []string{"ls", "--short"},
			Delete:    []string{"delete-session", "--force", "{session}"},
			Create:    []string{"attach", "--create-background", "{session}"},
			RenameTab: []string{"action", "rename-tab", "{title}"},
			Attach: Command{
				Pre: []string{"kitty", "@", "load-config", "{home}/.config/kitty/kitty-no-bind.conf"},
				Args: []string{
					"attach",
					"--create",
					"{session}",
				},
				Post: []string{"kitty", "@", "load-config", "{home}/.config/kitty/kitty.conf"},
			},
		},
	}
	ErrNotExists = fmt.Errorf("config not found: %w", os.ErrNotExist)
)

type Config struct {
	Zellij Zellij `yaml:"zellij"`
}

type Zellij struct {
	Bin       string   `yaml:"bin"`
	Ls        []string `yaml:"ls"`
	Delete    []string `yaml:"delete"`
	Create    []string `yaml:"create"`
	RenameTab []string `yaml:"rename_tab"`
	Attach    Command  `yaml:"attach"`
}

type Command struct {
	Pre  []string `yaml:"pre"`
	Args []string `yaml:"args"`
	Post []string `yaml:"post"`
}

func Default(filename string) error {
	return save(
		_default,
		filename,
	)
}

func New(filename string) (*Config, error) {
	var (
		config = _default
	)
	if _, err := os.Stat(filename); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &_default, ErrNotExists
		}
		return nil, fmt.Errorf("failed to find file: %w", err)
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &config, nil
}

func save(config any, path string) error {
	data, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("failed to marshall config: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o777); err != nil {
		return fmt.Errorf("failed to save default config: %w", err)
	}
	if err := os.WriteFile(path, data, 0o777); err != nil {
		return fmt.Errorf("failed to save default config: %w", err)
	}
	return nil
}
