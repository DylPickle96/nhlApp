package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

var nhlCities = []string{
	"Boston",
	"Buffalo",
	"Toronto",
	"Florida",
	"Montreal",
	"Tampa Bay",
	"Ottawa",
	"Detroit",
	"Washington",
	"NY Islanders",
	"Carolina",
	"Pittsburgh",
	"Philadelphia",
	"Columbus",
	"New Jersey",
	"NY Rangers",
	"St Louis",
	"Nashville",
	"Colorado",
	"Dallas",
	"Winnipeg",
	"Chicago",
	"Minnesota",
	"Edmonton",
	"Vancouver",
	"Arizona",
	"Vegas",
	"Calgary",
	"Anaheim",
	"San Jose",
	"Los Angeles",
}

type monthsRange struct {
	beginnning int
	ending     int
}

var twentyTwentySeason = map[string]map[string]monthsRange{
	"2020": {
		"Oct": {
			beginnning: 2,
			ending:     31,
		},
		"Nov": {
			beginnning: 1,
			ending:     19,
		},
		"Dec": {
			beginnning: 0,
			ending:     0,
		},
		"Jan": {
			beginnning: 0,
			ending:     0,
		},
		"Feb": {
			beginnning: 0,
			ending:     0,
		},
		"Mar": {
			beginnning: 0,
			ending:     0,
		},
		"Apr": {
			beginnning: 0,
			ending:     0,
		},
	},
}

type dailyRecord struct {
	TeamRecords []teamRecord `json:"teamRecords"`
}
type teamRecord struct {
	TeamName         string `json:"teamName"`
	Wins             string `json:"wins"`
	Loses            string `json:"loses"`
	Overtime         string `json:"overtime"`
	ROW              string `json:"ROW"`
	Points           string `json:"points"`
	GoalsFor         string `json:"goalsFor"`
	GoalsAgainst     string `json:"goalsAgainst"`
	Home             string `json:"home"`
	Away             string `json:"away"`
	DivisionRecord   string `json:"divisionRecord"`
	ConferenceRecord string `json:"conferenceRecord"`
	ICF              string `json:"ICF"`
}

func main() {
	err := getSeasonData()
	if err != nil {
		log.Printf("WARNING: %v", err)
	}
	// HTMLFile, err := os.Open("dylan.html")
	// if err != nil {
	// 	panic(err)
	// }
}

// fileParser main function to handle the parsing of the HTML file which are trying to scrape
func fileParser(HTMLFile []byte, season, month, day *string) {
	reader := bytes.NewReader(HTMLFile)
	tokenizer := html.NewTokenizer(reader)
	TDCount := 0
	var tr teamRecord
	var teamRecords []teamRecord
	var dailyRecord dailyRecord
	for {
		nextToken := tokenizer.Next()
		if nextToken == html.StartTagToken {
			token := tokenizer.Token()
			// webpage has the team name within an anchor tag when find one compare
			if token.Data == "a" {
				innerToken := tokenizer.Next()
				if innerToken == html.TextToken {
					value := (string)(tokenizer.Text())
					if !includesCity(value) {
						// this may produce bugs but the idea is if we do not find the city name we do not care about this block of HTML
						TDCount = 0
						continue
					}
					tr.TeamName = value
				}
			}
			// the table digit is the most important thing we care about this contains the data we are after
			if token.Data == "td" {
				innerToken := tokenizer.Next()
				var textValue string
				if innerToken == html.TextToken {
					textValue = strings.TrimSpace((string)(tokenizer.Text()))
				}
				tr = parsePostion(&TDCount, &textValue, tr)
				// append if this has been reset. Means we have reached the end of the row
				if TDCount == 10 {
					teamRecords = append(teamRecords, tr)
					TDCount = 0
					continue
				}
				// everytime we find a table digit we want to keep a count
				// for each team the things like wins, losses, etc will be found at the same count
				// this provides an easy interface to extract the data
				TDCount++
			}
		}
		// if we are at the end of the table row reset the count
		if nextToken == html.EndTagToken {
			token := tokenizer.Token()
			if token.Data == "tr" {
				TDCount = 0
			}
		}
		// if we are at the end of the file break
		if nextToken == html.ErrorToken {
			break
		}
	}
	teamRecords = validateData(teamRecords)
	// validate data function almost works the way it should. It ends up ignoring one of the two bad records at the beginnning of the slice
	// so the bad solution right now is call removeRecords here as well on the zeroth index
	// this poor practice but the more important thing here is that it works
	teamRecords = removeRecord(teamRecords, 0)
	dailyRecord.TeamRecords = teamRecords
	writeJSONFile(dailyRecord, season, month, day)
}

// parsePostion parse position assigns the value for that given position to the coorsponding teamRecord value
func parsePostion(TDCount *int, value *string, teamRecord teamRecord) teamRecord {
	switch count := *TDCount; count {
	case 1:
		winsLoses := strings.Split(*value, "-")
		if len(winsLoses) == 3 {
			teamRecord.Wins = winsLoses[0]
			teamRecord.Loses = winsLoses[1]
			teamRecord.Overtime = winsLoses[2]
		}
	case 2:
		teamRecord.ROW = *value
	case 3:
		teamRecord.Points = *value
	case 4:
		teamRecord.GoalsFor = *value
	case 5:
		teamRecord.GoalsAgainst = *value
	case 6:
		teamRecord.Home = *value
	case 7:
		teamRecord.Away = *value
	case 8:
		teamRecord.DivisionRecord = *value
	case 9:
		teamRecord.ConferenceRecord = *value
	case 10:
		teamRecord.ICF = *value
	}
	return teamRecord
}

// validateData checks if any of the values in a record are an empty string. This process can generate some blank records so we use this function to remove them
func validateData(teamRecords []teamRecord) []teamRecord {
	for i, tr := range teamRecords {
		if tr.Wins == "" {
			teamRecords = removeRecord(teamRecords, i)
			continue
		}
		if tr.Loses == "" {
			teamRecords = removeRecord(teamRecords, i)
			continue
		}
		if tr.Overtime == "" {
			teamRecords = removeRecord(teamRecords, i)
			continue
		}
		if tr.ROW == "" {
			teamRecords = removeRecord(teamRecords, i)
			continue
		}
		if tr.Points == "" {
			teamRecords = removeRecord(teamRecords, i)
			continue
		}
		if tr.GoalsFor == "" {
			teamRecords = removeRecord(teamRecords, i)
			continue
		}
		if tr.GoalsAgainst == "" {
			teamRecords = removeRecord(teamRecords, i)
			continue
		}
		if tr.Home == "" {
			teamRecords = removeRecord(teamRecords, i)
			continue
		}
		if tr.Away == "" {
			teamRecords = removeRecord(teamRecords, i)
			continue
		}
		if tr.DivisionRecord == "" {
			teamRecords = removeRecord(teamRecords, i)
			continue
		}
		if tr.ConferenceRecord == "" {
			teamRecords = removeRecord(teamRecords, i)
			continue
		}
		if tr.ICF == "" {
			teamRecords = removeRecord(teamRecords, i)
			continue
		}
	}
	return teamRecords
}

func getSeasonData() error {
	for season, months := range twentyTwentySeason {
		for month, monthRange := range months {
			if monthRange.beginnning == 0 {
				continue
			}
			for i := monthRange.beginnning; i <= monthRange.ending; i++ {
				resp, err := http.Get(fmt.Sprintf("http://www.shrpsports.com/nhl/stand.php?link=Y&season=%s&divcnf=div&month=%s&date=%d", season, month, i))
				if err != nil {
					return fmt.Errorf("could not get response from url %s. Error: %v", fmt.Sprintf("http://www.shrpsports.com/nhl/stand.php?link=Y&season=%s&divcnf=div&month=%s&date=%d", season, month, i), err)
				}
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return fmt.Errorf("could not read response body from url: %s. Error: %v", fmt.Sprintf("http://www.shrpsports.com/nhl/stand.php?link=Y&season=%s&divcnf=div&month=%s&date=%d", season, month, i), err)
				}
				day := strconv.FormatInt(int64(i), 10)
				fileParser(body, &season, &month, &day)
			}
		}
	}
	return nil
}

func writeJSONFile(dailyRecord dailyRecord, season, month, day *string) {
	JSONBytes, err := json.MarshalIndent(dailyRecord, "", "  ")
	if err != nil {
		log.Printf("WARNING: could not marshal slice of teamRecord into JSON. Error: %v", err)
	}
	err = ioutil.WriteFile(fmt.Sprintf("JSON/%s-%s-%s-record.json", *season, *month, *day), JSONBytes, 0755)
	if err != nil {
		log.Printf("WARNING: could not write JSON to file. Error: %v", err)
	}
}

// removeRecord removes a record from the slice of team records that are scraped from the website
func removeRecord(teamRecords []teamRecord, index int) []teamRecord {
	// fmt.Println("DEBUG: index: ", index)
	// fmt.Printf("DEBUG: teamRecord: %+v\n", teamRecords[index])
	teamRecords = append(teamRecords[:index], teamRecords[index+1:]...)
	return teamRecords
}

// includesCity checks if the passed in value is contained with the list of cities
func includesCity(value string) bool {
	for _, city := range nhlCities {
		if value == city {
			return true
		}
	}
	return false
}
