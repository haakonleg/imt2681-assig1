package apiinfo

import (
	"fmt"
	"net/http"
	"time"

	"github.com/haakonleg/imt2681-assig1/request"
)

// GetAPIInfo send information about the API
func GetAPIInfo(req *request.Request, startTime *time.Time) {
	response := struct {
		Uptime  string `json:"uptime"`
		Info    string `json:"info"`
		Version string `json:"version"`
	}{
		Uptime:  uptime(startTime),
		Info:    "Service for IGC tracks.",
		Version: "v1"}

	req.SendJSON(&response, http.StatusOK)
}

// uptime returns the app uptime in ISO 8601 duration format
func uptime(startTime *time.Time) string {
	// Seconds duration since app start
	duration := time.Since(*startTime)

	sec := int(duration.Seconds()) % 60
	min := int(duration.Minutes()) % 60
	hour := int(duration.Hours()) % 24
	day := int(duration.Hours()/24) % 7
	month := int(duration.Hours()/24/7/4) % 12
	year := int(duration.Hours() / 24 / 365)

	return fmt.Sprintf("P%dY%dM%dDT%dH%dM%dS", year, month, day, hour, min, sec)
}
