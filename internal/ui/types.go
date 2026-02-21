package ui

type refreshSuccessMsg string

type standingsMsg struct {
	currentMatchDay int
	teams           []team
}

type matchesMsg struct {
	matches []match
}

type teamMatchesMsg struct {
	matches []match
}

type selectedTeamMsg struct {
	teamName string
}

type errMsg struct {
	err error
}
