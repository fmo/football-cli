package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m model) refreshView() string {
	if m.refreshSuccess != "" {
		return m.refreshSuccess
	}

	return "do refresh"
}

func (m model) matchesTopView() string {
	commonStyle := lipgloss.NewStyle().PaddingLeft(2)

	buttonStyle := lipgloss.NewStyle().Width(40).Align(lipgloss.Right)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		commonStyle.Render(fmt.Sprintf("Matchday %d", m.currentMatchDay)),
		buttonStyle.Render("(p)rev | (n)ext"),
	)
}

func (m model) matchesView() string {
	bottomStyle := lipgloss.NewStyle()

	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.matchesTopView(),
		bottomStyle.Render(m.matchesTable.View()))
}

func (m model) RightView() string {
	baseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	if m.err != nil {
		return baseStyle.Render(m.err.Error())
	}

	if item, ok := m.list.SelectedItem().(item); ok {
		if item.Title() == "Matches" {
			return baseStyle.Render(m.matchesView())
		}
		if item.Title() == "Refresh Data" {
			return baseStyle.Render(m.refreshView())
		}
	}

	return baseStyle.Render(m.table.View())
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
