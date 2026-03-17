package view

const (
	joinSeparator = "  "

	appTitle      = "Zema 🧌"
	sessionsTitle = "sessions"
	helpTitle     = "help"

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

	statusCreatedPrefix        = "created:"
	statusDeletedPrefix        = "deleted:"
	statusCreatedMessagePrefix = statusCreatedPrefix + " "
	statusDeletedMessagePrefix = statusDeletedPrefix + " "
	statusSessionNotSelected   = "session is not selected ⚠️"
	statusSessionNameEmpty     = "session name is empty ⚠️"

	helpCloseHint = "Press any key to close help"
)
