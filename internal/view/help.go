package view

import "strings"

func (m model) viewHelp() string {
	return strings.Join([]string{
		join(titleStyle.Render("Zema:"), helpStyle.Render("help")),
		"",
		helpStyle.Render("j/k or up/down  move select"),
		helpStyle.Render("enter           attach to session"),
		helpStyle.Render("/               filter sessions"),
		helpStyle.Render("c               create session"),
		helpStyle.Render("d               delete session"),
		helpStyle.Render("q               quit"),
		"",
		helpStyle.Render("Press any key to close help"),
	}, "\n")
}
