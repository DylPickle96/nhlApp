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

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

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

// seasonNumericName required as I cannot use numeric value in a collection name...
var seasonNumericName = map[string]string{
	"2020": "twentyTwenty",
	"2019": "twentyNineteen",
}

type server struct {
	r *httprouter.Router
}

var client *mongo.Client

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
	handleRequests()
}

func getLeagueRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var dailyRecord dailyRecord
	filter := bson.D{{"month", ps.ByName("month")}, {"day", ps.ByName("day")}}
	collection := client.Database("nhlRecords").Collection(seasonNumericName[ps.ByName("season")] + "Season")
	err := collection.FindOne(context.Background(), filter).Decode(&dailyRecord)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("could not get record from database"))
		return
	}
	teamRecords, err := bubbleSort(dailyRecord.TeamRecords)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("could not sort records. Error: " + err.Error()))
		return
	}
	dailyRecord.TeamRecords = *teamRecords
	json.NewEncoder(w).Encode(dailyRecord)
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE")
	s.r.ServeHTTP(w, r)
}

func handleRequests() {
	myRouter := httprouter.New()
	myRouter.GET("/league/:season/:month/:day", getLeagueRecord)
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
