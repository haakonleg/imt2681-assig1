package igcinfo

import (
	"encoding/json"
	"fmt"
	"net/http"

	igc "github.com/marni/goigc"
)

type registerTrackRequest struct {
	URL string `json:"url"`
}

// In-memory storage of all tracks
var trackDbCount = 0
var trackDb = make(map[int]igc.Track)

// RegisterTrack registers a track in the database
func RegisterTrack(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request registerTrackRequest
	err := ParseJSONRequest(r, &request)
	if err != nil {
		WriteError(w, "Invalid JSON", 400)
		return
	}

	track, err := igc.ParseLocation(request.URL)
	if err != nil {
		fmt.Println(err.Error())
		WriteError(w, "Error parsing track location", 400)
		return
	}

	trackDbCount++
	trackDb[trackDbCount] = track
	WriteMessage(w, "Track registered successfully")
}

// GetAllTracks sends a response containing all registered track IDs
func GetAllTracks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ids := make([]int, 0, len(trackDb))
	for k := range trackDb {
		ids = append(ids, k)
	}

	res, _ := json.Marshal(ids)
	fmt.Fprintln(w, string(res))
}
