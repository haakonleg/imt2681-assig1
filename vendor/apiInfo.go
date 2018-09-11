package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AppStartTime is the start time of the app
var appStartTime = time.Now()

type apiInfo struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

// APIInfo sends information about the API
func APIInfo(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Content-Type", "application/json")

	info := &apiInfo{uptime(), "Service for IGC tracks.", "v1"}
	res, _ := json.Marshal(info)

	fmt.Fprintf(*w, string(res))
}

func uptime() string {
	duration := time.Since(appStartTime).Round(time.Second)

	fmt.Println(duration)
	return fmt.Sprintln(duration)
}
