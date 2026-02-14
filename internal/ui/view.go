package ui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m model) RightView() string {
	if m.err != nil {
		return baseStyle.Render(m.err.Error())
	}

	if item, ok := m.list.SelectedItem().(item); ok {
		if item.Title() == "Matches" {
			return baseStyle.Render(m.matchesTable.View())
		}
		if item.Title() == "Refresh Data" {
			if m.refreshSuccess != "" {
				return baseStyle.Render(m.refreshSuccess)
			} else {
				return baseStyle.Render("do refresh")
			}
		}
	}
	return baseStyle.Render(m.table.View())
}

func (m model) View() string {
	leftStyle := lipgloss.NewStyle().Width(40).Padding(0, 1)
	rightStyle := lipgloss.NewStyle().Width(80).Padding(0, 1)

	left := leftStyle.Render(m.list.View())
	right := rightStyle.Render(m.RightView())

	page := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	return docStyle.Render(page)
}
