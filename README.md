# ZEMA

`zema` is a terminal TUI tool written in Go for managing [Zellij](https://zellij.dev/) sessions.

## Features

- list Zellij sessions;
- filter sessions by search input;
- create a new session;
- delete the selected session;
- attach to a session directly from the interface;
- built-in help screen with keyboard shortcuts.

## Build

Recommended approach (via `Taskfile.yml`):

```bash
task build
```

The binary will be built to `bin/zema`.

## Controls

- `j/k` or `up/down` - move through the list;
- `enter` - attach to the selected session;
- `/` - filter sessions;
- `c` - create a session;
- `d` - delete a session;
- `?` - open help;
- `q` or `ctrl+c` - quit.
