package igcinfo

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/haakonleg/imt2681-assig1/request"
	"github.com/haakonleg/imt2681-assig1/track"
	"github.com/haakonleg/imt2681-assig1/apiinfo"
)

const apiPath = "/igcinfo/api/"

// App contains the context of the API and its required members/variables
// The App struct must be instantiated with a ListenPort set to set the listening port of the API
type App struct {
	ListenPort string

	db        map[int]track.Track
	startTime time.Time
}

// Route the API request to handlers
func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := request.CreateRequest(w, r, r.Method)
	path := strings.TrimPrefix(r.URL.Path, apiPath)

	// GET /api
	if len(path) == 0 && req.Method == request.GET {
		apiinfo.GetAPIInfo(req, &app.startTime)
		return
	}

	// This regex matches every possible API path by checking which optional capture groups are zero or non-zero (i.e they were matched or not)
	var reIGC = regexp.MustCompile("^igc/?(/[0-9]+)?/?(/(pilot|glider|glider_id|track_length|H_date|track_src_url))?/?$")
	if match := reIGC.FindStringSubmatch(path); match != nil {
		// Matches GET /api/igc and POST /api/igc
		if len(match[1]) == 0 && len(match[2]) == 0 {
			switch req.Method {
			case request.GET:
				track.GetAllTracks(req, &app.db)
			case request.POST:
				track.RegisterTrack(req, &app.db)
			}
			return
			// Matches GET /api/igc/{id}
		} else if len(match[1]) != 0 && len(match[2]) == 0 && req.Method == request.GET {
			// The id should always be a real int because of the regex, so can ignore error
			id, _ := strconv.Atoi(match[1][1:])
			track.GetTrack(req, &app.db, id)
			return
			// Matches GET /api/igc/{id}/{field}
		} else if len(match[1]) != 0 && len(match[2]) != 0 && req.Method == request.GET {
			id, _ := strconv.Atoi(match[1][1:])
			track.GetTrackField(req, &app.db, id, match[2][1:])
			return
		}
	}

	http.NotFound(w, r)
}

// StartServer starts the API HTTP server
func (app *App) StartServer() {
	app.startTime = time.Now()

	if app.ListenPort == "" {
		log.Fatal("ListenPort must be set")
	}

	app.db = make(map[int]Track, 0)

	// Add HTTP handler
	http.Handle(apiPath, app)

	// Start listen
	fmt.Printf("Server listening on port %s\n", app.ListenPort)
	if err := http.ListenAndServe(":"+app.ListenPort, nil); err != nil {
		log.Fatal(err.Error())
	}
}
