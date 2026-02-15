package ui

import (
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "n":
			if item, ok := m.list.SelectedItem().(item); ok {
				if item.Title() == "Matches" {
					m.currentMatchDay++
					return m, matchesHandler(m.currentMatchDay)
				}
			}
		case "p":
			if item, ok := m.list.SelectedItem().(item); ok {
				if item.Title() == "Matches" {
					m.currentMatchDay--
					return m, matchesHandler(m.currentMatchDay)
				}
			}
		case "enter":
			if item, ok := m.list.SelectedItem().(item); ok {
				if item.Title() == "Matches" {
					return m, matchesHandler(m.currentMatchDay)
				}
				if item.Title() == "Refresh Data" {
					return m, refreshHandler
				}
			}
		}
	case matchesMsg:
		m.matches = msg.matches
		matchesMap := make(map[string][]match)

		for _, match := range msg.matches {
			matchDate := match.utcDate.Format(time.DateOnly)
			matchesMap[matchDate] = append(matchesMap[matchDate], match)
		}

		for _, mm := range matchesMap {
			matchesTable := NewMatchesTable()
			matchesTable.table.SetRows(buildMatches(mm))
		}

		//m.matchesTable.SetRows(buildMatches(msg.matches))
	case standingsMsg:
		m.currentMatchDay = msg.currentMatchDay
		m.teams = msg.teams
		m.table.SetRows(buildRows(msg.teams))
		return m, matchesHandler(msg.currentMatchDay)
	case refreshSuccessMsg:
		m.refreshSuccess = string(msg)
	case errMsg:
		m.err = msg.err
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var listCmd, tableCmd, matchesTableCmd tea.Cmd

	m.list, listCmd = m.list.Update(msg)
	m.table, tableCmd = m.table.Update(msg)
	//m.matchesTable, matchesTableCmd = m.matchesTable.Update(msg)

	return m, tea.Batch(tableCmd, matchesTableCmd, listCmd)
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
