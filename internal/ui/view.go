package ui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m model) RightView() string {
	baseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	s := ""

	if m.err != nil {
		s += baseStyle.Render(m.err.Error())
	}

	if item, ok := m.list.SelectedItem().(item); ok {
		if item.Title() == "Matches" {
			s += baseStyle.Render(m.matchesTable.View())
		}
		if item.Title() == "Refresh Data" {
			if m.refreshSuccess != "" {
				s += baseStyle.Render(m.refreshSuccess)
			} else {
				s += baseStyle.Render("do refresh")
			}
		}
	}

	if s == "" {
		s += m.table.View()
	}

	topStyle := lipgloss.NewStyle()
	bottomStyle := lipgloss.NewStyle()

	page := lipgloss.JoinVertical(
		lipgloss.Top,
		topStyle.Render("26th Match Day"),
		bottomStyle.Render(s))

	return baseStyle.Render(page)
}

func (m model) View() string {
	leftStyle := lipgloss.NewStyle().Width(40).Padding(0, 1)
	rightStyle := lipgloss.NewStyle().Width(80).Padding(0, 1)

	page := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStyle.Render(m.list.View()),
		rightStyle.Render(m.RightView()))

	return docStyle.Render(page)
}
