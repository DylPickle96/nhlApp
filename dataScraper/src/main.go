package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/net/html"
)

type seasons struct {
	seasons []season
}

type season struct {
	name   string
	months []month
}

type month struct {
	name      string
	beginning int
	ending    int
}
type dailyRecord struct {
	Season      string       `json:"season"`
	Month       string       `json:"month"`
	Day         string       `json:"day"`
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

var (
	client    *mongo.Client
	nhlCities = []string{
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
		"Seattle",
	}
	nhlSeasons = seasons{
		seasons: []season{
			season{
				name: "2020",
				months: []month{
					month{
						name:      "Oct",
						beginning: 2,
						ending:    31,
					},
					month{
						name:      "Nov",
						beginning: 1,
						ending:    30,
					},
					month{
						name:      "Dec",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Jan",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Feb",
						beginning: 1,
						ending:    29,
					},
					month{
						name:      "Mar",
						beginning: 1,
						ending:    11,
					},
				},
			},
			season{
				name: "2019",
				months: []month{
					month{
						name:      "Oct",
						beginning: 3,
						ending:    31,
					},
					month{
						name:      "Nov",
						beginning: 1,
						ending:    30,
					},
					month{
						name:      "Dec",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Jan",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Feb",
						beginning: 2,
						ending:    28,
					},
					month{
						name:      "Mar",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Apr",
						beginning: 1,
						ending:    6,
					},
				},
			},
			season{
				name: "2018",
				months: []month{
					month{
						name:      "Oct",
						beginning: 4,
						ending:    31,
					},
					month{
						name:      "Nov",
						beginning: 1,
						ending:    30,
					},
					month{
						name:      "Dec",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Jan",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Feb",
						beginning: 2,
						ending:    28,
					},
					month{
						name:      "Mar",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Apr",
						beginning: 1,
						ending:    8,
					},
				},
			},
			season{
				name: "2017",
				months: []month{
					month{
						name:      "Oct",
						beginning: 12,
						ending:    31,
					},
					month{
						name:      "Nov",
						beginning: 1,
						ending:    30,
					},
					month{
						name:      "Dec",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Jan",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Feb",
						beginning: 2,
						ending:    28,
					},
					month{
						name:      "Mar",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Apr",
						beginning: 1,
						ending:    9,
					},
				},
			},
			season{
				name: "2016",
				months: []month{
					month{
						name:      "Oct",
						beginning: 7,
						ending:    31,
					},
					month{
						name:      "Nov",
						beginning: 1,
						ending:    30,
					},
					month{
						name:      "Dec",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Jan",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Feb",
						beginning: 2,
						ending:    28,
					},
					month{
						name:      "Mar",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Apr",
						beginning: 1,
						ending:    10,
					},
				},
			},
			season{
				name: "2015",
				months: []month{
					month{
						name:      "Oct",
						beginning: 8,
						ending:    31,
					},
					month{
						name:      "Nov",
						beginning: 1,
						ending:    30,
					},
					month{
						name:      "Dec",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Jan",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Feb",
						beginning: 2,
						ending:    28,
					},
					month{
						name:      "Mar",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Apr",
						beginning: 1,
						ending:    11,
					},
				},
			},
			season{
				name: "2014",
				months: []month{
					month{
						name:      "Oct",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Nov",
						beginning: 1,
						ending:    30,
					},
					month{
						name:      "Dec",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Jan",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Feb",
						beginning: 2,
						ending:    28,
					},
					month{
						name:      "Mar",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Apr",
						beginning: 1,
						ending:    13,
					},
				},
			},
			season{
				name: "2013",
				months: []month{
					month{
						name:      "Jan",
						beginning: 19,
						ending:    31,
					},
					month{
						name:      "Feb",
						beginning: 2,
						ending:    28,
					},
					month{
						name:      "Mar",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Apr",
						beginning: 1,
						ending:    28,
					},
				},
			},
			season{
				name: "2012",
				months: []month{
					month{
						name:      "Oct",
						beginning: 6,
						ending:    31,
					},
					month{
						name:      "Nov",
						beginning: 1,
						ending:    30,
					},
					month{
						name:      "Dec",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Jan",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Feb",
						beginning: 2,
						ending:    28,
					},
					month{
						name:      "Mar",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Apr",
						beginning: 1,
						ending:    7,
					},
				},
			},
			season{
				name: "2011",
				months: []month{
					month{
						name:      "Oct",
						beginning: 7,
						ending:    31,
					},
					month{
						name:      "Nov",
						beginning: 1,
						ending:    30,
					},
					month{
						name:      "Dec",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Jan",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Feb",
						beginning: 2,
						ending:    28,
					},
					month{
						name:      "Mar",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Apr",
						beginning: 1,
						ending:    10,
					},
				},
			},
			season{
				name: "2010",
				months: []month{
					month{
						name:      "Oct",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Nov",
						beginning: 1,
						ending:    30,
					},
					month{
						name:      "Dec",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Jan",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Feb",
						beginning: 2,
						ending:    28,
					},
					month{
						name:      "Mar",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Apr",
						beginning: 1,
						ending:    11,
					},
				},
			},
			season{
				name: "2009",
				months: []month{
					month{
						name:      "Oct",
						beginning: 4,
						ending:    31,
					},
					month{
						name:      "Nov",
						beginning: 1,
						ending:    30,
					},
					month{
						name:      "Dec",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Jan",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Feb",
						beginning: 2,
						ending:    28,
					},
					month{
						name:      "Mar",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Apr",
						beginning: 1,
						ending:    12,
					},
				},
			},
			season{
				name: "2008",
				months: []month{
					month{
						name:      "Sep",
						beginning: 29,
						ending:    30,
					},
					month{
						name:      "Oct",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Nov",
						beginning: 1,
						ending:    30,
					},
					month{
						name:      "Dec",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Jan",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Feb",
						beginning: 2,
						ending:    28,
					},
					month{
						name:      "Mar",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Apr",
						beginning: 1,
						ending:    6,
					},
				},
			},
			season{
				name: "2007",
				months: []month{
					month{
						name:      "Oct",
						beginning: 4,
						ending:    31,
					},
					month{
						name:      "Nov",
						beginning: 1,
						ending:    30,
					},
					month{
						name:      "Dec",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Jan",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Feb",
						beginning: 2,
						ending:    28,
					},
					month{
						name:      "Mar",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Apr",
						beginning: 1,
						ending:    11,
					},
				},
			},
			season{
				name: "2006",
				months: []month{
					month{
						name:      "Oct",
						beginning: 5,
						ending:    31,
					},
					month{
						name:      "Nov",
						beginning: 1,
						ending:    30,
					},
					month{
						name:      "Dec",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Jan",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Feb",
						beginning: 2,
						ending:    28,
					},
					month{
						name:      "Mar",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Apr",
						beginning: 1,
						ending:    18,
					},
				},
			},
			season{
				name: "2004",
				months: []month{
					month{
						name:      "Oct",
						beginning: 8,
						ending:    31,
					},
					month{
						name:      "Nov",
						beginning: 1,
						ending:    30,
					},
					month{
						name:      "Dec",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Jan",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Feb",
						beginning: 2,
						ending:    28,
					},
					month{
						name:      "Mar",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Apr",
						beginning: 1,
						ending:    4,
					},
				},
			},
			season{
				name: "2003",
				months: []month{
					month{
						name:      "Oct",
						beginning: 9,
						ending:    31,
					},
					month{
						name:      "Nov",
						beginning: 1,
						ending:    30,
					},
					month{
						name:      "Dec",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Jan",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Feb",
						beginning: 2,
						ending:    28,
					},
					month{
						name:      "Mar",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Apr",
						beginning: 1,
						ending:    6,
					},
				},
			},
			season{
				name: "2002",
				months: []month{
					month{
						name:      "Oct",
						beginning: 3,
						ending:    31,
					},
					month{
						name:      "Nov",
						beginning: 1,
						ending:    30,
					},
					month{
						name:      "Dec",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Jan",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Feb",
						beginning: 2,
						ending:    28,
					},
					month{
						name:      "Mar",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Apr",
						beginning: 1,
						ending:    14,
					},
				},
			},
			season{
				name: "2001",
				months: []month{
					month{
						name:      "Oct",
						beginning: 4,
						ending:    31,
					},
					month{
						name:      "Nov",
						beginning: 1,
						ending:    30,
					},
					month{
						name:      "Dec",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Jan",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Feb",
						beginning: 2,
						ending:    28,
					},
					month{
						name:      "Mar",
						beginning: 1,
						ending:    31,
					},
					month{
						name:      "Apr",
						beginning: 1,
						ending:    8,
					},
				},
			},
		},
	}
	// seasonNumericName - required as I cannot use numeric value in a Mongodb collection name...
	seasonNumericName = map[string]string{
		"2020": "twentyTwenty",
		"2019": "twentyNineteen",
		"2018": "twentyEighteen",
		"2017": "twentySeventeen",
		"2016": "twentySixteen",
		"2015": "twentyFifteen",
		"2014": "twentyFourteen",
		"2013": "twentyThirteen",
		"2012": "twentyTwelve",
		"2011": "twentyEleven",
		"2010": "twentyTen",
		"2009": "twoThousandAndNine",
		"2008": "twoThousandAndEight",
		"2007": "twoThousandAndSeven",
		"2006": "twoThousandAndSix",
		"2005": "twoThousandAndFive",
		"2004": "twoThousandAndFour",
		"2003": "twoThousandAndThree",
		"2002": "twoThousandAndTwo",
		"2001": "twoThousandAndOne",
	}
)

func init() {
	var err error
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("ERROR: Could not connect to the mongo database. Error: %v", err)
	}
	ctx, cancelFunc = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("ERROR: could not ping the mongo database. Error: %v", err)
	}
}

func main() {
	err := getSeasonsData()
	if err != nil {
		log.Printf("WARNING: %v", err)
	}
}

// fileParser - main function to handle the parsing of the HTML file which are trying to scrape
func fileParser(HTMLFile []byte, season, month, day string) {
	reader := bytes.NewReader(HTMLFile)
	tokenizer := html.NewTokenizer(reader)
	TDCount := 0
	// season numeric
	s, _ := strconv.ParseInt(season, 10, 64)
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
				tr = parsePostion(TDCount, s, textValue, tr)
				// append if this has been reset. Means we have reached the end of the row
				if s <= 2010 && TDCount == 9 {
					teamRecords = append(teamRecords, tr)
					TDCount = 0
					continue
				} else if TDCount == 10 {
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
	fmt.Println(season, month, day)
	if len(teamRecords) == 1 {
		return
	}
	teamRecords = validateData(teamRecords, s)
	// validate data function almost works the way it should. It ends up ignoring one of the two bad records at the beginnning of the slice
	// so the bad solution right now is call removeRecords here as well on the zeroth index
	// this poor practice but the more important thing here is that it works
	teamRecords = removeRecord(teamRecords, 0)
	dailyRecord.TeamRecords = teamRecords
	dailyRecord.Season = season
	dailyRecord.Month = month
	dailyRecord.Day = day
	insertDailyRecord(dailyRecord, season)
	writeJSONFile(dailyRecord, season, month, day)
}

// parsePostion - parse position assigns the value for that given position to the coorsponding teamRecord value
func parsePostion(TDCount int, season int64, value string, teamRecord teamRecord) teamRecord {
	if season <= 2010 {
		switch count := TDCount; count {
		case 1:
			winsLoses := strings.Split(value, "-")
			if len(winsLoses) == 4 && season <= 2004 {
				teamRecord.Wins = winsLoses[0]
				teamRecord.Loses = winsLoses[1]
				teamRecord.Overtime = winsLoses[2] + winsLoses[3]
			} else if len(winsLoses) == 3 {
				teamRecord.Wins = winsLoses[0]
				teamRecord.Loses = winsLoses[1]
				teamRecord.Overtime = winsLoses[2]
			}
		case 2:
			teamRecord.Points = value
		case 3:
			teamRecord.GoalsFor = value
		case 4:
			teamRecord.GoalsAgainst = value
		case 5:
			teamRecord.Home = value
		case 6:
			teamRecord.Away = value
		case 7:
			teamRecord.DivisionRecord = value
		case 8:
			teamRecord.ConferenceRecord = value
		case 9:
			teamRecord.ICF = value
		}
	} else {
		switch count := TDCount; count {
		case 1:
			winsLoses := strings.Split(value, "-")
			if len(winsLoses) == 3 {
				teamRecord.Wins = winsLoses[0]
				teamRecord.Loses = winsLoses[1]
				teamRecord.Overtime = winsLoses[2]
			}
		case 2:
			teamRecord.ROW = value
		case 3:
			teamRecord.Points = value
		case 4:
			teamRecord.GoalsFor = value
		case 5:
			teamRecord.GoalsAgainst = value
		case 6:
			teamRecord.Home = value
		case 7:
			teamRecord.Away = value
		case 8:
			teamRecord.DivisionRecord = value
		case 9:
			teamRecord.ConferenceRecord = value
		case 10:
			teamRecord.ICF = value
		}
	}

	return teamRecord
}

// validateData - checks if any of the values in a record are an empty string. This process can generate some blank records so we use this function to remove them
func validateData(teamRecords []teamRecord, season int64) []teamRecord {
	for i, tr := range teamRecords {
		if tr.TeamName == "" {
			teamRecords = removeRecord(teamRecords, i)
			continue
		}
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
		if season > 2010 && tr.ROW == "" {
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

// getSeasonsData - gets season data from the shrpsports website
func getSeasonsData() error {
	for _, season := range nhlSeasons.seasons {
		s, _ := strconv.ParseInt(season.name, 10, 64)
		if s != 2004 {
			continue
		}
		for _, month := range season.months {
			for i := month.beginning; i <= month.ending; i++ {
				resp, err := http.Get(fmt.Sprintf("http://www.shrpsports.com/nhl/stand.php?link=Y&season=%s&divcnf=div&month=%s&date=%d", season.name, month.name, i))
				if err != nil {
					return fmt.Errorf("could not get response from url %s. Error: %v", fmt.Sprintf("http://www.shrpsports.com/nhl/stand.php?link=Y&season=%s&divcnf=div&month=%s&date=%d", season.name, month.name, i), err)
				}
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return fmt.Errorf("could not read response body from url: %s. Error: %v", fmt.Sprintf("http://www.shrpsports.com/nhl/stand.php?link=Y&season=%s&divcnf=div&month=%s&date=%d", season.name, month.name, i), err)
				}
				day := strconv.FormatInt(int64(i), 10)
				fileParser(body, season.name, month.name, day)
			}
		}
	}
	return nil
}

// insertDailyRecord - insert daily records into the mongoDB collection
func insertDailyRecord(dailyRecord dailyRecord, season string) {
	collection := client.Database("nhlRecords").Collection(seasonNumericName[season] + "Season")
	insertRecord, err := collection.InsertOne(context.Background(), dailyRecord)
	if err != nil {
		log.Printf("WARNING: could not insert record in collection %s-season.Error: %v", season, err)
		return
	}
	log.Println("insertRecord ID: ", insertRecord.InsertedID)
}

func writeJSONFile(dailyRecord dailyRecord, season, month, day string) {
	JSONBytes, err := json.MarshalIndent(dailyRecord, "", "  ")
	if err != nil {
		log.Printf("WARNING: could not marshal slice of teamRecord into JSON. Error: %v", err)
	}
	err = os.MkdirAll(fmt.Sprintf("JSON/%s", season), 0755)
	if err != nil {
		log.Fatalf("ERROR: could not make the directory path JSON/%s. Error: %v", season, err)
	}
	err = ioutil.WriteFile(fmt.Sprintf("JSON/%s/%s-%s-%s-record.json", season, season, month, day), JSONBytes, 0755)
	if err != nil {
		log.Printf("WARNING: could not write JSON to file. Error: %v", err)
	}
}

// removeRecord - removes a record from the slice of team records that are scraped from the website
func removeRecord(teamRecords []teamRecord, index int) []teamRecord {
	teamRecords = append(teamRecords[:index], teamRecords[index+1:]...)
	return teamRecords
}

// includesCity - checks if the passed in value is contained with the list of cities
func includesCity(value string) bool {
	for _, city := range nhlCities {
		if value == city {
			return true
		}
	}
	return false
}
