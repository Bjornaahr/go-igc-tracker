package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/marni/goigc"
)

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

func handlerF(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "F to pay respect")
}

//GetPort retrives the port
func GetPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "4747"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}

func main() {
	http.HandleFunc("/p/", handlerPilot)
	http.HandleFunc("/", handlerF)
	http.ListenAndServe(GetPort(), nil)
}
