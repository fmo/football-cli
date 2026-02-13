package ui

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	footballdataapi "github.com/fmo/football-data-api"
)

func refreshHandler() tea.Msg {
	if err := os.Remove("data/standings.json"); err != nil {
		log.Println("cant refresh data: ", err)
	}
	return refreshSuccessMsg("refreshed the data")
}

func standingsHandler() tea.Msg {
	// check if there is standings in the data folder
	f, err := os.Open("data/standings.json")
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	var resp *footballdataapi.RespCompStandings

	// could not open file so read it from api
	if err != nil {
		log.Println("making api call")
		client := footballdataapi.NewClient(&http.Client{Timeout: 10 * time.Second})
		compReq := footballdataapi.NewReqCompStandings(client)

		resp, err = compReq.Do()
		if err != nil {
			return nil
		}

		f, err = os.OpenFile("data/standings.json", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Println("opening file error: ", err)
		}
		defer f.Close()

		jsonStandings, err := json.Marshal(resp)
		if err != nil {
			log.Println("marshalling error: ", err)
		}

		_, err = f.Write(jsonStandings)
		if err != nil {
			log.Println("file write error: ", err)
		}
	} else {
		log.Println("reading data from the file")

		r, err := io.ReadAll(f)
		if err != nil {
			log.Println("cant read the data from the json file: ", err)
		}
		if err := json.Unmarshal(r, &resp); err != nil {
			log.Println("cant unmarshal: ", err)
			return errMsg{err}
		}
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
