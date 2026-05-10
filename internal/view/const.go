package view

const (
	joinSeparator = "  "

	keyQuit   = "q"
	keyCtrlC  = "ctrl+c"
	keyHelp   = "?"
	keyCreate = "c"
	keyEnter  = "enter"
	keyDelete = "d"
	keyEsc    = "esc"

	actionQuit       = "quit"
	actionHelp       = "help"
	actionMoveSelect = "move select"
	actionAttach     = "attach to session"
	actionFilter     = "filter sessions"
	actionCreate     = "create session"
	actionDelete     = "delete session"

	labelNewSession    = "New session"
	labelNoSessions    = "No sessions"
	labelChooseSession = "Choose session"
	descTotalFmt       = "Total: %d"

	statusCreatedPrefix      = "created:"
	statusDeletedPrefix      = "deleted:"
	statusAttachedPrefix     = "attached:"
	statusSessionNotSelected = "session is not selected"
	statusSessionNameEmpty   = "session name is empty"
)
