package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/marni/goigc"
	"github.com/p3lim/iso8601"
)

const (
	VERSION = "1.0"
	DESC    = "Service for IGC tracks."
)

//MetaInfo about the program
type MetaInfo struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

//Glider track info
type Track struct {
	ID          int       `json: "ID"`
	Hdate       time.Time `json:"H_date"`
	Pilot       string    `json:"pilot"`
	Glider      string    `json:"glider"`
	GliderID    string    `json:"glider_id"`
	Tracklength string    `json:"track_length"`
}

var startTime time.Time
var tracks map[int]Track
var ID int

func init() {
	startTime = time.Now()
	tracks = make(map[int]Track)
	ID = 1
}

func handlerPilot(w http.ResponseWriter, r *http.Request) {
	s := "http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"
	track, err := igc.ParseLocation(s)
	if err != nil {
		status := 400
		http.Error(w, http.StatusText(status), status)
		return
	}

	fmt.Fprintln(w, track.Pilot)
}

func handlerAPI(w http.ResponseWriter, r *http.Request) {

	info := MetaInfo{iso8601.Format(time.Since(startTime)),
		DESC,
		VERSION}
	infoJSON, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(infoJSON)

}

func handlerIGC(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case ("GET"):
		IDS := []int{}

		for id := range tracks {
			IDS = append(IDS, id)
		}
		IDJSON, err := json.Marshal(IDS)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(IDJSON)

	//Creates a track
	case ("POST"):
		var url string
		err := json.NewDecoder(r.Body).Decode(&url)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		track, err := igc.ParseLocation(url)
		tracks[ID] = Track{ID, track.Date, track.Pilot, track.GliderType, track.GliderID, "Length"}

		infoJSON, err := json.Marshal(tracks[ID].ID)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(infoJSON)
		ID++
	}
}

//GetPort retrives the port from the enviorment
func GetPort() string {
	//Gets the port
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "5000"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}

func main() {
	http.HandleFunc("/api/", handlerAPI)
	http.HandleFunc("/api/igc", handlerIGC)
	http.ListenAndServe(GetPort(), nil)
}
