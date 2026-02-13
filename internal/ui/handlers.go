package ui

import (
	"fmt"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	footballdataapi "github.com/fmo/football-data-api"
)

func standingsHandler() tea.Msg {
	client := footballdataapi.NewClient(&http.Client{Timeout: 10 * time.Second})
	compReq := footballdataapi.NewReqCompStandings(client)

	resp, err := compReq.Do()
	if err != nil {
		return nil
	}

	sm := standingsMsg{}
	for _, s := range resp.Standings {
		team := team{}
		for _, t := range s.Table {
			team.position = t.Position
			team.name = t.Team.Name
			team.points = t.Points
			sm.teams = append(sm.teams, team)
		}
	}

	return sm
}

func matchesHandler() tea.Msg {
	client := footballdataapi.NewClient(&http.Client{Timeout: 10 * time.Second})
	request := footballdataapi.NewMatches(client)

	resp, err := request.Do(footballdataapi.PL, footballdataapi.Finished, 2025, 25)
	if err != nil {
		return nil
	}

	mm := matchesMsg{}
	for _, m := range resp.Matches {
		match := match{}
		match.homeTeam = m.HomeTeam.Name
		match.awayTeam = m.AwayTeam.Name
		match.score = fmt.Sprintf("%d-%d", m.Score.FullTime.Home, m.Score.FullTime.Away)
		mm.matches = append(mm.matches, match)
	}

	return mm
}
