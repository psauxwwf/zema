package view

import "strings"

func (m model) viewHelp() string {
	return strings.Join([]string{
		titleStyle.Render("help"),
		"",
		helpStyle.Render("j/k or up/down  move"),
		helpStyle.Render("/               filter"),
		helpStyle.Render("a               add item"),
		helpStyle.Render("q               quit"),
		"",
		helpStyle.Render("Press any key to close help"),
	}, "\n")
}
