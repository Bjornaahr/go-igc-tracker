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

//MetaInfo about the program
type MetaInfo struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

var startTime time.Time

func init() {
	startTime = time.Now()
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

	switch r.Method {

	case ("GET"):
		info := MetaInfo{iso8601.Format(time.Since(startTime)),
			"Service for IGC tracks",
			"1.0"}
		infoJSON, err := json.Marshal(info)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(infoJSON)

	case ("POST"):

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
	http.HandleFunc("/API/", handlerAPI)
	http.ListenAndServe(GetPort(), nil)
}
