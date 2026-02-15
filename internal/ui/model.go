package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
)

type model struct {
	list            list.Model
	teams           []team
	table           table.Model
	matchesTable    table.Model
	matches         []match
	refreshSuccess  string
	err             error
	currentMatchDay int
}

type team struct {
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
