package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

type teamRecord struct {
	teamName         string
	wins             string
	loses            string
	overtime         string
	ROW              string
	points           string
	goalsFor         string
	goalsAgainst     string
	home             string
	away             string
	divisionRecord   string
	conferenceRecord string
	icf              string
}

func main() {
	resp, err := http.Get("http://www.shrpsports.com/nhl/stand.php?link=Y&season=2020&divcnf=div&month=Nov&date=5")
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// HTMLFile, err := os.Open("dylan.html")
	// if err != nil {
	// 	panic(err)
	// }
	fileParser(body)
}

func fileParser(HTMLFile []byte) {
	reader := bytes.NewReader(HTMLFile)
	tokenizer := html.NewTokenizer(reader)
	TDCount := 0
	var tr teamRecord
	var teamRecords []teamRecord
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
					tr.teamName = value
					fmt.Println("City: ", value)
				}
			}
			// the table digit is the most important thing we care about this contains the data we are after
			if token.Data == "td" {
				innerToken := tokenizer.Next()
				var textValue string
				if innerToken == html.TextToken {
					textValue = strings.TrimSpace((string)(tokenizer.Text()))
					fmt.Printf("count: %d. Value: %s\n", TDCount, textValue)
				}
				tr = parsePostion(&TDCount, &textValue, tr)
				// append if this has been reset. Means we have reached the
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
	count := 0
	for _, teamRecord := range teamRecords {
		if includesCity(teamRecord.teamName) {
			log.Printf("teamRecord: %+v\n", teamRecord)
			count++
		}
	}
	log.Println("count: ", count)
}

func parsePostion(TDCount *int, value *string, teamRecord teamRecord) teamRecord {
	switch count := *TDCount; count {
	case 1:
		winsLoses := strings.Split(*value, "-")
		if len(winsLoses) == 3 {
			teamRecord.wins = winsLoses[0]
			teamRecord.loses = winsLoses[1]
			teamRecord.overtime = winsLoses[2]
		}
	case 2:
		teamRecord.ROW = *value
	case 4:
		teamRecord.points = *value
	case 5:
		teamRecord.goalsFor = *value
	case 6:
		teamRecord.goalsAgainst = *value
	case 7:
		teamRecord.home = *value
	case 8:
		teamRecord.away = *value
	case 9:
		teamRecord.divisionRecord = *value
	case 10:
		teamRecord.icf = *value
	}
	return teamRecord
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
