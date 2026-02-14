package ui

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	footballdataapi "github.com/fmo/football-data-api"
)

const (
	standingsFile = "data/standings.json"
	matchesFile   = "data/matches.json"
)

const seasonYear = 2025

func refreshHandler() tea.Msg {
	if err := os.Remove(standingsFile); err != nil {
		log.Println("cant refresh data: ", err)
	}

	if err := os.Remove(matchesFile); err != nil {
		log.Println("cant refresh match data")
	}

	return refreshSuccessMsg("refreshed the data")
}

func standingsHandler() tea.Msg {
	fileExist, err := fileExist(standingsFile)
	if err != nil {
		log.Println(err)
	}

	var resp *footballdataapi.RespCompStandings

	// could not open file so read it from api
	if !fileExist {
		log.Println("making api call for standings")

		client := footballdataapi.NewClient(&http.Client{Timeout: 10 * time.Second})
		compReq := footballdataapi.NewReqCompStandings(client)

		resp, err = compReq.Do(footballdataapi.PL, seasonYear)
		if err != nil {
			return nil
		}

		data, err := json.Marshal(resp)
		if err != nil {
			log.Println("marshalling error: ", err)
		}

		if err := writeData(data, standingsFile); err != nil {
			log.Println("writing data problem")
		}
	} else {
		log.Println("reading standing data from the file")

		r, err := readData(standingsFile)
		if err != nil {
			log.Println(err)
		}

		if err := json.Unmarshal(r, &resp); err != nil {
			log.Println("cant unmarshal: ", err)
			return errMsg{errors.New("cant unmarshal the data")}
		}
	}

	sm := standingsMsg{}
	sm.currentMatchDay = resp.Season.CurrentMatchDay
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

func matchesHandler(currentMatchDay int) tea.Cmd {
	return func() tea.Msg {
		matchesFileExist, err := fileExist(matchesFile)
		if err != nil {
			log.Println("cant check if match data file exist: ", err)
		}

		if currentMatchDay == 0 {
			return nil
		}

		var resp *footballdataapi.RespMatches

		if !matchesFileExist {
			log.Printf("requesting matches for season: %d, matchday: %d\n", seasonYear, currentMatchDay)

			client := footballdataapi.NewClient(&http.Client{Timeout: 10 * time.Second})
			request := footballdataapi.NewMatches(client)

			resp, err = request.Do(footballdataapi.PL, seasonYear, currentMatchDay)
			if err != nil {
				return nil
			}

			data, err := json.Marshal(resp)
			if err != nil {
				log.Println("marshaling error: ", err)
			}

			if err := writeData(data, matchesFile); err != nil {
				log.Println("writing data problem")
			}
		} else {
			log.Println("reading matches data form the file")

			r, err := readData(matchesFile)
			if err != nil {
				log.Println(err)
			}

			if err := json.Unmarshal(r, &resp); err != nil {
				log.Println("cant unmarshal: ", err)
				return errMsg{errors.New("cant unmarshal the data")}
			}
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
}

func readData(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r, err := io.ReadAll(f)
	if err != nil {
		log.Println("cant read the data from the json file: ", err)
	}

	return r, nil
}

func fileExist(filename string) (bool, error) {
	f, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer f.Close()

	return true, nil
}

func writeData(data []byte, filename string) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println("opening file error: ", err)
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		log.Println("file write error: ", err)
		return err
	}

	return nil
}
