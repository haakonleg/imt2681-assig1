package track

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/haakonleg/imt2681-assig1/request"
	igc "github.com/marni/goigc"
)

// Track is the model of IGC tracks stored in database
type Track struct {
	HDate       string `json:"H_Date"`
	Pilot       string `json:"pilot"`
	Glider      string `json:"glider"`
	GliderID    string `json:"glider_id"`
	TrackLength string `json:"track_length"`
}

// Creates a new track object out of a parsed IGC track from goigc
func createTrack(igc *igc.Track) Track {
	return Track{
		HDate:       igc.Date.String(),
		Pilot:       igc.Pilot,
		Glider:      igc.GliderType,
		GliderID:    igc.GliderID,
		TrackLength: calTrackLen(&igc.Points).String()}
}

// Calculate track length, time of last point subtracted by time of first
func calTrackLen(points *[]igc.Point) time.Duration {
	arrLen := len(*points)

	firstTime := (*points)[0].Time
	lastTime := (*points)[arrLen-1].Time

	return lastTime.Sub(firstTime)
}

// Ensures that a link points to an IGC resource (but just that it is a valid URL and has an igc extension)
func ensureIGCLink(link string) bool {
	if _, err := url.ParseRequestURI(link); err != nil {
		return false
	}

	ext := strings.ToLower(path.Ext(link))
	if ext != ".igc" {
		return false
	}
	return true
}

// GET /api/track
// Returns an array of IDs of all tracks stored in the database
func GetAllTracks(req *request.Request, db *map[int]Track) {
	// Get all track IDs in database
	ids := make([]int, 0, len(*db))
	for id := range *db {
		ids = append(ids, id)
	}

	req.SendJSON(&ids, http.StatusOK)
}

// POST /api/track
// Register/upload a track
func RegisterTrack(req *request.Request, db *map[int]Track) {
	var request struct {
		URL string `json:"url"`
	}

	// Get the JSON post request
	if err := req.ParseJSONRequest(&request); err != nil {
		req.SendError("Error parsing JSON request", http.StatusBadRequest)
		return
	}

	// Check that the supplied link is valid
	if valid := ensureIGCLink(request.URL); !valid {
		req.SendError("This is not a valid IGC link", http.StatusBadRequest)
		return
	}

	// Parse the IGC file
	igc, err := igc.ParseLocation(request.URL)
	if err != nil {
		fmt.Println(err)
		req.SendError("Error parsing IGC track", http.StatusBadRequest)
		return
	}

	// Create a new track and store it in the database
	id := len(*db)
	(*db)[id] = createTrack(&igc)

	// Send response containing the ID to the inserted track
	response := struct {
		ID int `json:"id"`
	}{ID: id}
	req.SendJSON(&response, http.StatusOK)
}

// GET /api/track/{id}
// Retrieves a track by the value of its ObjectID (hex encoded string)
func GetTrack(req *request.Request, db *map[int]Track, id int) {
	// Check if track is in the database
	track, ok := (*db)[id]
	if !ok {
		req.SendError("Invalid ID", http.StatusBadRequest)
		return
	}

	req.SendJSON(&track, http.StatusOK)
}

// GET /api/track/{id}/{field}
// Retrieves a field in the track object
func GetTrackField(req *request.Request, db *map[int]Track, id int, field string) {
	// Check if track is in the database
	track, ok := (*db)[id]
	if !ok {
		req.SendError("Invalid ID", http.StatusBadRequest)
		return
	}

	switch field {
	case "pilot":
		req.SendText(track.Pilot)
	case "glider":
		req.SendText(track.Glider)
	case "glider_id":
		req.SendText(track.GliderID)
	case "track_length":
		req.SendText(track.TrackLength)
	case "H_date":
		req.SendText(track.HDate)
	}
}
