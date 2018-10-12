package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
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
	Hdate       time.Time `json:"date from File Header, H-record"`
	Pilot       string    `json:"pilot"`
	Glider      string    `json:"glider"`
	GliderID    string    `json:"glider_id"`
	Tracklength float64   `json:"calculated total track length"`
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

func handlerGetTrack(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	varID := vars["id"]
	id, err := strconv.Atoi(varID)
	if err != nil {
		panic(err)
	}
	TrackJSON, err := json.Marshal(tracks[id])
	if err != nil {
		panic(err)
	}
	if tracks[id].ID == 0 {
		http.Error(w, "404 Not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(TrackJSON)
}

func handlerGetField(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	varID := vars["id"]
	id, err := strconv.Atoi(varID)
	if err != nil {
		panic(err)
	}

	trackMap := map[string]string{
		"pilot":        tracks[id].Pilot,
		"glider":       tracks[id].Glider,
		"glider_id":    tracks[id].GliderID,
		"track_length": fmt.Sprintf("%f", tracks[id].Tracklength),
		"h_date":       tracks[id].Hdate.String(),
	}
	field := strings.ToLower(vars["field"])

	if val, ok := trackMap[field]; ok {
		TrackJSON, err := json.Marshal(val)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(TrackJSON)
	} else {
		http.Error(w, "400 - Bad Request, the field you entered is not on our database!", http.StatusBadRequest)
		return

	}

}

func handlerIGC(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case ("GET"):

		ids := []int{}

		for index := range tracks {
			ids = append(ids, tracks[index].ID)
		}
		IDJSON, err := json.Marshal(ids)
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
		tracks[ID] = Track{ID, track.Date, track.Pilot, track.GliderType, track.GliderID, CalculateDistance(track)}

		infoJSON, err := json.Marshal(tracks[ID].ID)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(infoJSON)
		ID++

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

//CalculateDistance calculates the track distance
func CalculateDistance(track igc.Track) float64 {

	trackdistance := 0.0
	//Loops through all the points and find the distance between them
	for i := 0; i < len(track.Points)-1; i++ {
		trackdistance += track.Points[i].Distance(track.Points[i+1])
	}

	return trackdistance
}

//GetPort retrives the port from the enviorment
func GetPort() string {
	//Gets the port
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "5000"
		fmt.Println("Could not find port in enviorment, setting port to: " + port)
	}
	return ":" + port
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/igcinfo/api/", handlerAPI)
	router.HandleFunc("/igcinfo/api/igc/", handlerIGC)
	router.HandleFunc("/igcinfo/api/igc/{id:[0-9]+}/", handlerGetTrack)
	router.HandleFunc("/igcinfo/api/igc/{id:[0-9]+}/{field:[a-zA-Z_]+}/", handlerGetField)
	http.ListenAndServe(GetPort(), router)
}
