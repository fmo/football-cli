package ui

type refreshSuccessMsg string

type standingsMsg struct {
	currentMatchDay int
	teams           []team
}

type matchesMsg struct {
	matches []match
}

type errMsg struct {
	err error
}
