package main

import (
	"fmt"
	"net/http"

	"github.com/marni/goigc"
)

func pilotHandler(w http.ResponseWriter, r *http.Request) {
	s := "http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"
	track, err := igc.ParseLocation(s)
	if err != nil {
		status := 400
		http.Error(w, http.StatusText(status), status)
		return
	}

	fmt.Fprintf(w, "Pilot: %s, gliderType: %s, date: %s",
		track.Pilot, track.GliderType, track.Date.String())
}

func main() {

}
