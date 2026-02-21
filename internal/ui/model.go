package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
)

type model struct {
	list               list.Model
	teams              []team
	table              table.Model
	matchesTables      []table.Model
	matchesDates       []string
	matches            []match
	teamMatches        []match
	refreshSuccess     string
	err                error
	currentMatchDay    int
	tableSelectionView bool
	teamDetailView     bool
	selectedTeam       string
}

type team struct {
	id       int
	position int
	name     string
	points   int
}

type match struct {
	homeTeam string
	awayTeam string
	score    string
	utcDate  time.Time
}
