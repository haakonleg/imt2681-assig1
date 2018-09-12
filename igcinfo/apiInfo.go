package igcinfo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var appStartTime = time.Now()

type apiInfo struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

// APIInfo sends information about the API
func APIInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	info := &apiInfo{uptime(), "Service for IGC tracks.", "v1"}
	res, _ := json.Marshal(info)

	fmt.Fprintf(w, string(res))
}

// uptime returns the app uptime in ISO 8601 duration format
func uptime() string {
	// Seconds duration since app start
	duration := time.Since(appStartTime)

	sec := int(duration.Seconds()) % 60
	min := int(duration.Minutes()) % 60
	hour := int(duration.Hours()) % 24
	day := int(duration.Hours()/24) % 7
	month := int(duration.Hours()/24/7/4) % 12
	year := int(duration.Hours() / 24 / 365)

	return fmt.Sprintf("P%dY%dM%dDT%dH%dM%dS", year, month, day, hour, min, sec)
}
