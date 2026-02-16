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
		// update matches model with retrieved matches from the api
		m.matches = msg.matches

		// need a map for the matches
		matchesMap := make(map[string][]match)

		// need a map for match dates
		matchesDates := make([]string, 0)

		// put the matches into the map with the matchDate key 2026-02-10
		for _, match := range msg.matches {
			matchDate := match.utcDate.Format(time.DateOnly)
			if _, exists := matchesMap[matchDate]; !exists {
				matchesDates = append(matchesDates, matchDate)
			}

			matchesMap[matchDate] = append(matchesMap[matchDate], match)
		}

		m.matchesDates = matchesDates

		// create tables for each date
		for _, mm := range matchesMap {
			matchesTable := NewMatchesTable()
			matchesTable.table.SetRows(buildMatches(mm))
			m.matchesTables = append(m.matchesTables, matchesTable.table)
		}
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

	var listCmd, tableCmd tea.Cmd

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
