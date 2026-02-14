package ui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

type standingsTable struct {
	table table.Model
}

func (st standingsTable) New() *standingsTable {
	columns := []table.Column{
		{Title: "Position", Width: 8},
		{Title: "Team", Width: 30},
		{Title: "Points", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithHeight(21),
		table.WithWidth(50),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return &standingsTable{
		table: t,
	}
}

type matchesTable struct {
	table table.Model
}

func (mt matchesTable) New() *matchesTable {
	columns := []table.Column{
		{Title: "Home Team", Width: 20},
		{Title: "Away Team", Width: 20},
		{Title: "Score", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithHeight(12),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return &matchesTable{
		table: t,
	}
}
