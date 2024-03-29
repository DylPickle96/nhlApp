package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"go.mongodb.org/mongo-driver/mongo"
)

type seasonRecord struct {
	Season       string        `json:"season"`
	DailyRecords []dailyRecord `json:"dailyRecords"`
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

// http server interface
type server struct {
	r *httprouter.Router
}

var (
	client     *mongo.Client
	nhlSeasons = seasons{
		seasons: []season{
			season{
				name: "2021",
				months: []month{
					month{
						name:      "Jan",
						beginning: 13,
						ending:    31,
					},
					month{
						name:      "Feb",
						beginning: 1,
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
						ending:    30,
					},
					month{
						name:      "May",
						beginning: 1,
						ending:    18,
					},
				},
			},
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
		"2021": "twentyTwentyOne",
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
		log.Fatalf("ERROR: init(): Could not connect to the mongo database. Error: %v", err)
	}
	ctx, cancelFunc = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("ERROR: init(): could not ping the mongo database. Error: %v", err)
	}
}

func main() {
	handleRequests()
}

// getDailyLeagueRecord - obtains the entire leagues record for a daily in a season
func getDailyLeagueRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	dailyRecord, err := retrieveDailyLeagueRecord(ps.ByName("season"), ps.ByName("month"), ps.ByName("day"))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("could not sort records. Error: " + err.Error()))
		return
	}
	json.NewEncoder(w).Encode(dailyRecord)
}

// getSeasonRecord - HTTP interface to return a whole season worth of daily records
func getSeasonRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sr := seasonRecord{}
	sr.Season = ps.ByName("season")
	for _, season := range nhlSeasons.seasons {
		if season.name == ps.ByName("season") {
			for _, month := range season.months {
				for i := month.beginning; i <= month.ending; i++ {
					r, err := retrieveDailyLeagueRecord(season.name, month.name, strconv.FormatInt(int64(i), 10))
					if err != nil {
						log.Printf("ERROR: getSeasonRecord(): cannot retrieve daily record for season: %s month: %s. day: %d. Error: %v", season.name, month.name, i, err)
						w.WriteHeader(500)
						return
					}
					sr.DailyRecords = append(sr.DailyRecords, *r)
				}
			}
			break
		}
	}
	err := json.NewEncoder(w).Encode(sr)
	if err != nil {
		log.Printf("ERROR: getSeasonRecord(): cannot encode response. Error: %v", err)
		w.WriteHeader(500)
		return
	}
}

// retrieveDailyLeagueRecord - retrieves daily league record from our mongodb database
func retrieveDailyLeagueRecord(season, month, day string) (*dailyRecord, error) {
	dailyRecord := &dailyRecord{}
	filter := bson.D{{"month", month}, {"day", day}}
	collection := client.Database("nhlRecords").Collection(seasonNumericName[season] + "Season")
	err := collection.FindOne(context.Background(), filter).Decode(&dailyRecord)
	if err != nil {
		log.Printf("ERROR: retrieveDailyLeagueRecord(): Cannot get daily league record. Error: %v", err)
		return nil, err
	}
	teamRecords, err := bubbleSort(dailyRecord.TeamRecords)
	if err != nil {
		log.Printf("ERROR: retrieveDailyLeagueRecord(): issue during bubble sort. Error: %v", err)
		return nil, err
	}
	dailyRecord.TeamRecords = *teamRecords
	return dailyRecord, nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	s.r.ServeHTTP(w, r)
}

func handleRequests() {
	myRouter := httprouter.New()
	myRouter.GET("/daily/:season/:month/:day", getDailyLeagueRecord)
	myRouter.GET("/season/:season", getSeasonRecord)
	myRouter.ServeFiles("/webapp/*filepath", http.Dir("../../webApp/dist/"))
	log.Println("INFO: Started http listener")
	log.Fatal(http.ListenAndServe(":8081", &server{myRouter}))
}

func bubbleSort(teamRecords []teamRecord) (*[]teamRecord, error) {
	sorting := true
	for sorting {
		sorting = false
		for i := 1; i < len(teamRecords); i++ {
			previousValue, err := strconv.ParseInt(teamRecords[i-1].Points, 10, 64)
			if err != nil {
				return nil, err
			}
			currentValue, err := strconv.ParseInt(teamRecords[i].Points, 10, 64)
			if err != nil {
				return nil, err
			}
			// if the previous value is less than the current value swap it. This will sort the slice in descending order.
			if previousValue < currentValue {
				teamRecords[i], teamRecords[i-1] = teamRecords[i-1], teamRecords[i]
				sorting = true
			}
		}
	}
	return &teamRecords, nil
}
