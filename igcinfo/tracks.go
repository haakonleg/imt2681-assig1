package igcinfo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	igc "github.com/marni/goigc"
)

type registerTrackRequest struct {
	URL string `json:"url"`
}

type getTrackResponse struct {
	Hdate       string `json:"H_date"`
	Pilot       string `json:"pilot"`
	Glider      string `json:"glider"`
	GliderID    string `json:"glider_id"`
	TrackLength string `json:"track_length"`
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

	if len(request.URL) == 0 {
		WriteError(w, "URL not provided", 400)
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
	fmt.Fprint(w, string(res))
}

// GetTrackByID sends a response containing meta information about a registered track
func GetTrackByID(w http.ResponseWriter, r *http.Request, id int) {
	w.Header().Set("Content-Type", "application/json")

	track, ok := trackDb[id]
	if !ok {
		http.NotFound(w, r)
		return
	}

	trackInfo := &getTrackResponse{
		track.Date.String(),
		track.Pilot,
		track.GliderType,
		track.GliderID,
		calTrackTime(&track.Points).String()}

	res, _ := json.Marshal(trackInfo)
	fmt.Fprint(w, string(res))
}

// GetTrackField sends a response containing a single field in a registered track
func GetTrackField(w http.ResponseWriter, r *http.Request, id int, field string) {
	w.Header().Set("Content-Type", "text/plain")

	track, ok := trackDb[id]
	if !ok {
		http.NotFound(w, r)
		return
	}

	switch field {
	case "pilot":
		fmt.Fprint(w, track.Pilot)
	case "glider":
		fmt.Fprint(w, track.GliderType)
	case "glider_id":
		fmt.Fprint(w, track.GliderID)
	case "track_length":
		fmt.Fprint(w, calTrackTime(&track.Points).String())
	case "H_date":
		fmt.Fprint(w, track.Date.String())
	default:
		http.NotFound(w, r)
	}
}

// Calculate track length, time of last point subtracted by time of first
func calTrackTime(points *[]igc.Point) time.Duration {
	arrLen := len(*points)

	firstTime := (*points)[0].Time
	lastTime := (*points)[arrLen-1].Time

	return lastTime.Sub(firstTime)
}
