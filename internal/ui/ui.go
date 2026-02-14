// Package ui renders the UI
package ui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().MarginTop(10).MarginLeft(10)

func (m model) Init() tea.Cmd {
	return standingsHandler
}

func Render() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Generate menu and standings table
	m := model{
		list:         menu{}.New().list,
		table:        standingsTable{}.New().table,
		matchesTable: matchesTable{}.New().table,
	}

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
