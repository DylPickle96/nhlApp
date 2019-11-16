package main

import (
	"fmt"
	"os"
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
	HTMLFile, err := os.Open("dylan.html")
	if err != nil {
		panic(err)
	}
	fileParser(HTMLFile)
}

func fileParser(HTMLFile *os.File) {
	tokenizer := html.NewTokenizer(HTMLFile)
	TDCount := 0
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
					fmt.Println("city: ", value)
				}
			}
			// the table digit is the most important thing we care about this contains the data we are after
			if token.Data == "td" {
				innerToken := tokenizer.Next()
				if innerToken == html.TextToken {
					value := (string)(tokenizer.Text())
					fmt.Printf("count: %d. Value: %s\n", TDCount, strings.TrimSpace(value))
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
