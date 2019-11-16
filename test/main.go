package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type nhlRecords struct {
	Copyright string `json:"copyright"`
	Records   []struct {
		StandingsType string `json:"standingsType"`
		League        struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Link string `json:"link"`
		} `json:"league"`
		Division struct {
			ID           int    `json:"id"`
			Name         string `json:"name"`
			NameShort    string `json:"nameShort"`
			Link         string `json:"link"`
			Abbreviation string `json:"abbreviation"`
		} `json:"division"`
		Conference struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Link string `json:"link"`
		} `json:"conference"`
		TeamRecords []struct {
			Team struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				Link string `json:"link"`
			} `json:"team"`
			LeagueRecord struct {
				Wins   int    `json:"wins"`
				Losses int    `json:"losses"`
				Ot     int    `json:"ot"`
				Type   string `json:"type"`
			} `json:"leagueRecord"`
			GoalsAgainst       int    `json:"goalsAgainst"`
			GoalsScored        int    `json:"goalsScored"`
			Points             int    `json:"points"`
			DivisionRank       string `json:"divisionRank"`
			DivisionL10Rank    string `json:"divisionL10Rank"`
			DivisionRoadRank   string `json:"divisionRoadRank"`
			DivisionHomeRank   string `json:"divisionHomeRank"`
			ConferenceRank     string `json:"conferenceRank"`
			ConferenceL10Rank  string `json:"conferenceL10Rank"`
			ConferenceRoadRank string `json:"conferenceRoadRank"`
			ConferenceHomeRank string `json:"conferenceHomeRank"`
			LeagueRank         string `json:"leagueRank"`
			LeagueL10Rank      string `json:"leagueL10Rank"`
			LeagueRoadRank     string `json:"leagueRoadRank"`
			LeagueHomeRank     string `json:"leagueHomeRank"`
			WildCardRank       string `json:"wildCardRank"`
			Row                int    `json:"row"`
			GamesPlayed        int    `json:"gamesPlayed"`
			Streak             struct {
				StreakType   string `json:"streakType"`
				StreakNumber int    `json:"streakNumber"`
				StreakCode   string `json:"streakCode"`
			} `json:"streak"`
			LastUpdated time.Time `json:"lastUpdated"`
		} `json:"teamRecords"`
	} `json:"records"`
}

type leagueRecord struct {
	teamID     int
	teamName   string
	leagueRank string
}

func main() {
	var (
		nhl           nhlRecords
		leagueRecords []leagueRecord
	)

	resp, err := http.Get("https://statsapi.web.nhl.com/api/v1/standings?date=2019-10-10")
	if err != nil {
		log.Fatalf("Could not fetch NHL standings. Error: %v", err)
	}
	json.NewDecoder(resp.Body).Decode(&nhl)

	for _, record := range nhl.Records {
		for _, teamRecord := range record.TeamRecords {
			leagueRecord := leagueRecord{
				teamID:     teamRecord.Team.ID,
				teamName:   teamRecord.Team.Name,
				leagueRank: teamRecord.LeagueRank,
			}
			leagueRecords = append(leagueRecords, leagueRecord)
		}
	}
	sortedRecords, err := bubbleSort(leagueRecords)
	if err != nil {
		log.Fatalf("could not sort league rankings. Error: %v", err)
	}
	for _, leagueRecord := range *sortedRecords {
		fmt.Printf("%+v\n", leagueRecord)
	}
}

func bubbleSort(leagueRecords []leagueRecord) (*[]leagueRecord, error) {
	sorting := true
	for sorting {
		sorting = false
		for i := 1; i < len(leagueRecords); i++ {
			previousValue, err := strconv.ParseInt(leagueRecords[i-1].leagueRank, 10, 64)
			if err != nil {
				return nil, err
			}
			currentValue, err := strconv.ParseInt(leagueRecords[i].leagueRank, 10, 64)
			if err != nil {
				return nil, err
			}
			if previousValue > currentValue {
				leagueRecords[i], leagueRecords[i-1] = leagueRecords[i-1], leagueRecords[i]
				sorting = true
			}
		}
	}
	return &leagueRecords, nil
}
