package ui

import (
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.tableSelectionView {
		return m.UpdateTeamView(msg)
	}

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
				if item.Title() == "Standings" {
					m.tableSelectionView = true
				}
			}
		}
	case matchesMsg:
		m.matches = msg.matches
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

	var listCmd tea.Cmd

	m.list, listCmd = m.list.Update(msg)

	return m, listCmd
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

func (m model) UpdateTeamView(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			if m.teamDetailView {
				m.teamDetailView = false
				return m, nil
			}

			m.tableSelectionView = false
			m.table.SetCursor(0)
			return m, nil
		case "enter":
			m.teamDetailView = true
			return m, func() tea.Msg {
				return selectedTeamMsg{m.table.SelectedRow()[1]}
			}
		}
	case selectedTeamMsg:
		m.selectedTeam = msg.teamName
	}

	var tableCmd tea.Cmd

	m.table, tableCmd = m.table.Update(msg)

	return m, tableCmd
}
