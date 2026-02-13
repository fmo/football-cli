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
		case "enter":
			if it, ok := m.list.SelectedItem().(item); ok {
				if it.Title() == "Games" {
					return m, matchesHandler
				}
			}
		case "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case standingsMsg:
		m.teams = msg.teams
		m.table.SetRows(buildRows(msg.teams))
	case matchesMsg:
		m.matches = msg.matches
		m.table.SetRows(buildMatches(msg.matches))
	}

	var tableCmd tea.Cmd
	var listCmd tea.Cmd

	m.list, listCmd = m.list.Update(msg)
	m.table, tableCmd = m.table.Update(msg)

	return m, tea.Batch(tableCmd, listCmd)
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
