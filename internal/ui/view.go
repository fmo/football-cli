package ui

import (
	"fmt"
	"strings"

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

	allTables := []string{}
	i := 0
	for _, md := range m.matchesDates {
		if i == 0 {
			allTables = append(allTables, lipgloss.JoinVertical(
				lipgloss.Top,
				m.matchesTopView(),
				lipgloss.NewStyle().Padding(0, 2).BorderStyle(lipgloss.NormalBorder()).Render(md),
				bottomStyle.Render(m.matchesTables[i].View())),
			)
			i++
			continue
		}

		m.matchesTables[i].Columns()
		remaining := lipgloss.JoinVertical(
			lipgloss.Top,
			lipgloss.NewStyle().Padding(0, 2).BorderStyle(lipgloss.NormalBorder()).Render(md),
			bottomStyle.Render(m.matchesTables[i].View()),
		)

		allTables = append(allTables, remaining)

		i++
	}

	return strings.Join(allTables, "\n")
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
	rightStyle := lipgloss.NewStyle().Width(100).Padding(0, 1)

	page := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStyle.Render(m.list.View()),
		rightStyle.Render(m.RightView()))

	return docStyle.Render(page)
}
