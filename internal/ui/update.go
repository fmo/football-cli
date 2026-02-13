package ui

import (
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "enter":
			if item, ok := m.list.SelectedItem().(item); ok {
				if item.Title() == "Games" {
					return m, matchesHandler
				}
			}
		}
	case matchesMsg:
		m.matches = msg.matches
		m.matchesTable.SetRows(buildMatches(msg.matches))
	case standingsMsg:
		m.teams = msg.teams
		m.table.SetRows(buildRows(msg.teams))
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var tableCmd tea.Cmd
	var listCmd tea.Cmd
	var matchesTableCmd tea.Cmd

	m.list, listCmd = m.list.Update(msg)
	m.table, tableCmd = m.table.Update(msg)
	m.matchesTable, matchesTableCmd = m.matchesTable.Update(msg)

	return m, tea.Batch(tableCmd, listCmd, matchesTableCmd)
}

func buildRows(teams []team) []table.Row {
	rows := []table.Row{}
	for _, team := range teams {
		position := strconv.Itoa(team.position)
		points := strconv.Itoa(team.points)
		rows = append(rows, table.Row{position, team.name, points})
	}

	return rows
}

func buildMatches(matches []match) []table.Row {
	rows := []table.Row{}
	for _, match := range matches {
		rows = append(rows, table.Row{match.homeTeam, match.awayTeam, match.score})
	}
	return rows
}
