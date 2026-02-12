package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	footballdataapi "github.com/fmo/football-data-api"
)

var docStyle = lipgloss.NewStyle().MarginTop(10).MarginLeft(10)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type item struct {
	title       string
	description string
}

func (i item) Title() string {
	return i.title
}

func (i item) Description() string {
	return i.description
}

func (i item) FilterValue() string {
	return i.title
}

type model struct {
	list  list.Model
	teams []team
	table table.Model
}

type team struct {
	position int
	name     string
	points   int
}

type tableMsg struct {
	teams []team
	table table.Model
}

func (m *model) standingsMsg() tea.Msg {
	client := footballdataapi.NewClient(&http.Client{Timeout: 10 * time.Second})
	compReq := footballdataapi.NewReqCompStandings(client)

	resp, err := compReq.Do()
	if err != nil {
		return nil
	}

	tm := tableMsg{}
	for _, s := range resp.Standings {
		team := team{}
		for _, t := range s.Table {
			team.position = t.Position
			team.name = t.Team.Name
			team.points = t.Points
			tm.teams = append(tm.teams, team)
		}
	}

	rows := []table.Row{}

	for _, t := range m.teams {
		position := strconv.Itoa(t.position)
		points := strconv.Itoa(t.points)
		rows = append(rows, table.Row{position, t.name, points})
	}

	m.table.SetRows(rows)

	return tm
}

func (m model) Init() tea.Cmd {
	return m.standingsMsg
}

func (m *model) RightView() string {
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case tableMsg:
		m.teams = msg.teams
	}

	var tableCmd tea.Cmd
	var listCmd tea.Cmd

	m.list, listCmd = m.list.Update(msg)
	m.table, tableCmd = m.table.Update(msg)

	return m, tea.Batch(tableCmd, listCmd)
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	columns := []table.Column{
		{Title: "Position", Width: 3},
		{Title: "Team", Width: 30},
		{Title: "Points", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(21),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := model{
		list: list.New([]list.Item{
			item{"Games", "Game List"},
			item{"Standings", "Standing"},
		}, list.NewDefaultDelegate(), 0, 0),
		table: t,
	}
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
