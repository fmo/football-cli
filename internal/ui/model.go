package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
)

type model struct {
	list    list.Model
	teams   []team
	table   table.Model
	matches []match
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
}
