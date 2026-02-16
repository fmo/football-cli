package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (m model) refreshView() string {
	if m.refreshSuccess != "" {
		return m.refreshSuccess
	}

	return "do refresh"
}

func (m model) matchesTopView() string {
	commonStyle := lipgloss.NewStyle()

	buttonStyle := lipgloss.NewStyle().Width(40).Align(lipgloss.Right)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		commonStyle.Render(fmt.Sprintf("Matchday %d", m.currentMatchDay)),
		buttonStyle.Render("(p)rev | (n)ext"),
	)
}

func (m model) matchesView() string {
	matchesMap := make(map[string][]match)

	matchesDates := make([]string, 0)

	for _, match := range m.matches {
		matchDate := match.utcDate.Format(time.DateOnly)
		if _, exists := matchesMap[matchDate]; !exists {
			matchesDates = append(matchesDates, matchDate)
		}

		matchesMap[matchDate] = append(matchesMap[matchDate], match)
	}

	dateStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("25")).
		Padding(0, 1).
		Bold(true)

	homeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Background(lipgloss.Color("236")).
		Padding(0, 1)

	scoreStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("31")).
		Padding(0, 1).
		Bold(true)

	rowStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Background(lipgloss.Color("237")).
		Padding(0, 1)

	maxHomeWidth := 0
	maxAwayWidth := 0
	for _, match := range m.matches {
		if w := lipgloss.Width(match.homeTeam); w > maxHomeWidth {
			maxHomeWidth = w
		}
		if w := lipgloss.Width(match.awayTeam); w > maxAwayWidth {
			maxAwayWidth = w
		}
	}

	// Keep the rows compact while still aligning team names.
	if maxHomeWidth > 24 {
		maxHomeWidth = 24
	}
	if maxAwayWidth > 24 {
		maxAwayWidth = 24
	}
	if maxHomeWidth < 12 {
		maxHomeWidth = 12
	}
	if maxAwayWidth < 12 {
		maxAwayWidth = 12
	}

	s := ""
	for _, md := range matchesDates {
		s += dateStyle.Render(md) + "\n"

		for _, match := range matchesMap[md] {
			kickoff := match.utcDate.Local().Format("15:04")
			home := lipgloss.NewStyle().Width(maxHomeWidth).Render(match.homeTeam)
			away := lipgloss.NewStyle().Width(maxAwayWidth).Render(match.awayTeam)

			row := lipgloss.JoinHorizontal(
				lipgloss.Top,
				homeStyle.Render(home),
				rowStyle.Render("vs"),
				homeStyle.Render(away),
				scoreStyle.Render(match.score),
				rowStyle.Render(kickoff),
			)

			s += row + "\n"
		}
		s += "\n"
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.matchesTopView(),
		"",
		s,
	)
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
