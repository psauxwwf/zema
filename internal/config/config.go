package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
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
	Bin       string   `yaml:"bin" env-default:"zellij"`
	Ls        []string `yaml:"ls" env-default:"ls,--short"`
	Delete    []string `yaml:"delete" env-default:"delete-session,--force,{session}"`
	Create    []string `yaml:"create" env-default:"attach,--create-background,{session}"`
	RenameTab []string `yaml:"rename_tab" env-default:"action,rename-tab,{title}"`
	Attach    Command  `yaml:"attach"`
}

type Command struct {
	Pre  []string `yaml:"pre" env-default:"kitty,@,load-config,{home}/.config/kitty/kitty-no-bind.conf"`
	Args []string `yaml:"args" env-default:"attach,--create,{session}"`
	Post []string `yaml:"post" env-default:"kitty,@,load-config,{home}/.config/kitty/kitty.conf"`
}

func Default(filename string) error {
	return save(
		_default,
		filename,
	)
}

func New(filename string) (*Config, error) {
	var (
		config Config
	)
	if _, err := os.Stat(filename); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &_default, ErrNotExists
		}
		return nil, fmt.Errorf("failed to find file: %w", err)
	}
	if err := cleanenv.ReadConfig(filename, &config); err != nil {
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
